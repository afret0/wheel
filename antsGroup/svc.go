// Package antsGroup 提供一个基于 ants 协程池、API 与 golang.org/x/sync/errgroup 对齐的并发任务组。
//
// 与 errgroup 的差异：
//  1. 并发上限由池大小决定（New/WithContext 时指定），而不是 SetLimit。
//  2. 任务内的 panic 不会导致进程崩溃，而是被捕获、打日志并转换成 *PanicError 由 Wait 返回。
//  3. Wait 之后该 Group 即失效（池已释放），不可复用；此时再调用 Go 会记录 ErrGroupClosed 而不是阻塞。
//
// 与 errgroup 相同的语义：
//   - 第一个非 nil error 会被记录并 cancel 掉 WithContext 返回的 ctx。
//   - Wait 返回第一个非 nil error。
//   - Go 不可与 Wait 并发调用（与 sync.WaitGroup 的约束一致）。
package antsGroup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/panjf2000/ants/v2"

	"github.com/afret0/wheel/log"
	"github.com/afret0/wheel/tool"
)

// defaultPoolSize 未显式指定且未配置 ANTS_POOL_SIZE 时的默认池大小
const defaultPoolSize = 50

var (
	// ErrGroupClosed Wait() 之后再调用 Go() 会记录该错误（任务不会被执行）
	ErrGroupClosed = errors.New("antsGroup: group already closed by Wait()")

	// ErrPoolUnavailable 协程池创建失败，Group 不可用
	ErrPoolUnavailable = errors.New("antsGroup: pool unavailable")
)

// PanicError 任务内 panic 被捕获后转换成的 error
type PanicError struct {
	Value any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("antsGroup: task panic: %v\n%s", e.Value, e.Stack)
}

// Unwrap 若 panic 的值本身就是 error，则支持 errors.Is / errors.As 穿透
func (e *PanicError) Unwrap() error {
	if err, ok := e.Value.(error); ok {
		return err
	}
	return nil
}

type Group struct {
	ctx    context.Context
	cancel context.CancelFunc

	pool *ants.Pool

	wg sync.WaitGroup

	once   sync.Once
	closed atomic.Bool

	err error
	mu  sync.Mutex
}

// WithContext 返回一个 Group 和一个派生 ctx。
// 第一个返回非 nil error（或发生 panic）的任务会 cancel 掉该 ctx；Wait 返回时也会 cancel。
func WithContext(ctx context.Context, poolSize int) (*Group, context.Context) {
	if poolSize <= 0 {
		poolSize = defaultPoolSize
	}

	c, cancel := context.WithCancel(ctx)
	g := &Group{
		ctx:    c,
		cancel: cancel,
	}

	// 兜底 handler：正常路径下 panic 已在 g.run 内被 recover，
	// 走到这里说明是 ants 内部或 wrapper 之外的异常，至少要留下日志。
	pool, err := ants.NewPool(poolSize, ants.WithPanicHandler(func(p any) {
		log.GetLogger().Errorf("antsGroup: pool level panic: %v\n%s", p, debug.Stack())
	}))
	if err != nil {
		log.GetLogger().Errorf("antsGroup: create pool failed, size: %d, err: %s", poolSize, err)
		g.err = fmt.Errorf("%w: %s", ErrPoolUnavailable, err)
		g.closed.Store(true)
		cancel()
		return g, c
	}

	g.pool = pool
	return g, c
}

// New 创建一个使用 context.Background() 的 Group。
// 池大小优先级：显式入参 > 环境变量 ANTS_POOL_SIZE > defaultPoolSize。
func New(poolSizeChain ...int) *Group {
	poolSize := defaultPoolSize
	if PS := tool.ConStringToInt64WithoutErr(os.Getenv("ANTS_POOL_SIZE")); PS > 0 {
		poolSize = int(PS)
	}
	for _, v := range poolSizeChain {
		if v > 0 {
			poolSize = v
		}
	}

	g, _ := WithContext(context.Background(), poolSize)
	return g
}

// Go 提交一个任务。池满时会阻塞等待空闲 worker（与 errgroup.SetLimit 的行为一致）。
// 不可与 Wait 并发调用。
func (g *Group) Go(f func() error) {
	if f == nil {
		return
	}

	// Wait 已经把池释放了，此时再提交任务只会拿到 ErrPoolClosed，
	// 提前拦截避免 wg 计数与实际执行不一致。
	if g.closed.Load() {
		g.setErr(ErrGroupClosed)
		return
	}

	if g.pool == nil {
		g.setErr(ErrPoolUnavailable)
		return
	}

	g.wg.Add(1)

	if err := g.pool.Submit(func() { g.run(f) }); err != nil {
		// 关键：Submit 失败时 wrapper 不会执行，必须在这里补 Done，
		// 否则 wg 计数永远归不了零，Wait() 将永久阻塞。
		g.wg.Done()
		log.GetLogger().Errorf("antsGroup: submit failed: %s", err)
		g.setErr(err)
	}
}

// run 执行单个任务，负责 panic 捕获、ctx 短路与 error 收集。
func (g *Group) run(f func() error) {
	defer g.wg.Done()

	// ants 的 worker 是复用的，任务 panic 绝不能逃逸出去，
	// 否则会连带影响后续复用该 worker 的任务。
	defer func() {
		if r := recover(); r != nil {
			pe := &PanicError{Value: r, Stack: debug.Stack()}
			log.GetLogger().Errorf("%s", pe.Error())
			g.setErr(pe)
		}
	}()

	// 已经有任务失败或上游取消了，后续任务直接跳过
	select {
	case <-g.ctx.Done():
		return
	default:
	}

	if err := f(); err != nil {
		g.setErr(err)
	}
}

// setErr 只记录第一个 error，并触发 cancel
func (g *Group) setErr(err error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.err == nil {
		g.err = err
		g.cancel()
	}
}

// Wait 等待所有已提交任务结束，释放协程池并返回第一个 error。
// 多次调用是安全的，返回值一致；但 Wait 之后 Group 不可再复用。
func (g *Group) Wait() error {
	g.once.Do(func() {
		g.closed.Store(true)
		g.wg.Wait()
		if g.pool != nil {
			g.pool.Release()
		}
		// 与 errgroup 一致：Wait 返回时释放 ctx，避免 context 泄漏
		g.cancel()
	})

	g.mu.Lock()
	defer g.mu.Unlock()
	return g.err
}

// Running 返回当前正在执行任务的 worker 数，用于监控埋点
func (g *Group) Running() int {
	if g.pool == nil {
		return 0
	}
	return g.pool.Running()
}

// Free 返回当前空闲的 worker 数，用于监控埋点
func (g *Group) Free() int {
	if g.pool == nil {
		return 0
	}
	return g.pool.Free()
}

