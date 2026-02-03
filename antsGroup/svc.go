package antsGroup

import (
	"context"
	"os"
	"sync"

	"github.com/panjf2000/ants/v2"

	"github.com/afret0/wheel/tool"
)

type Group struct {
	ctx    context.Context
	cancel context.CancelFunc

	pool *ants.Pool

	wg sync.WaitGroup

	once sync.Once
	err  error
	mu   sync.Mutex
}

func WithContext(ctx context.Context, poolSize int) (*Group, context.Context) {
	c, cancel := context.WithCancel(ctx)

	pool, _ := ants.NewPool(poolSize, ants.WithPanicHandler(func(p interface{}) {
		// panic recover：把 panic 当成 error 返回
	}))

	return &Group{
		ctx:    c,
		cancel: cancel,
		pool:   pool,
	}, c
}

func New(poolSizeChain ...int) *Group {
	poolSize := 50
	if PS := tool.ConStringToInt64WithoutErr(os.Getenv("ANTS_POOL_SIZE")); PS > 0 {
		poolSize = int(PS)
	}
	for _, v := range poolSizeChain {
		poolSize = v
	}

	g, _ := WithContext(context.Background(), poolSize)
	return g
}

func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	// 提交到 ants 池
	_ = g.pool.Submit(func() {
		defer g.wg.Done()

		// 如果 context 已取消，直接跳过
		select {
		case <-g.ctx.Done():
			return
		default:
		}

		// 执行任务
		if err := f(); err != nil {
			g.mu.Lock()
			if g.err == nil {
				g.err = err
				g.cancel() // 第一个错误触发取消
			}
			g.mu.Unlock()
		}
	})
}

func (g *Group) Wait() error {
	g.wg.Wait()
	g.pool.Release()
	return g.err
}
