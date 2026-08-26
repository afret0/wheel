package antsGroup

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// waitWithTimeout 在超时时间内等待 Wait() 返回，返回 (err, 是否超时)
func waitWithTimeout(g *Group, d time.Duration) (error, bool) {
	type result struct{ err error }
	ch := make(chan result, 1)
	go func() { ch <- result{g.Wait()} }()

	select {
	case r := <-ch:
		return r.err, false
	case <-time.After(d):
		return nil, true
	}
}

// TestPanicConvertedToError 回归：任务 panic 必须被捕获并转换成 error，而不是被静默吞掉
func TestPanicConvertedToError(t *testing.T) {
	g := New(4)

	var otherDone atomic.Bool
	g.Go(func() error { panic("boom") })
	g.Go(func() error {
		time.Sleep(10 * time.Millisecond)
		otherDone.Store(true)
		return nil
	})

	err, timeout := waitWithTimeout(g, 3*time.Second)
	if timeout {
		t.Fatal("Wait() 超时，panic 处理导致死锁")
	}
	if err == nil {
		t.Fatal("panic 应该被转换成 error，实际返回 nil")
	}

	var pe *PanicError
	if !errors.As(err, &pe) {
		t.Fatalf("期望 *PanicError，实际 %T: %v", err, err)
	}
	if pe.Value != "boom" {
		t.Errorf("期望 panic value = boom, 实际 %v", pe.Value)
	}
	if len(pe.Stack) == 0 {
		t.Error("PanicError 应该带上堆栈")
	}
	if !otherDone.Load() {
		t.Error("其他任务不应该被 panic 影响")
	}
}

// TestPanicWithErrorValue panic 值本身是 error 时，errors.Is 应该能穿透
func TestPanicWithErrorValue(t *testing.T) {
	sentinel := errors.New("sentinel")

	g := New(2)
	g.Go(func() error { panic(sentinel) })

	err, timeout := waitWithTimeout(g, 3*time.Second)
	if timeout {
		t.Fatal("Wait() 超时")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("errors.Is 应该能穿透 PanicError，实际 %v", err)
	}
}

// TestPanicDoesNotBreakPool 回归：panic 不能污染被复用的 worker
func TestPanicDoesNotBreakPool(t *testing.T) {
	g := New(1) // 池大小 1，强制复用同一个 worker

	var done atomic.Int32
	for i := 0; i < 5; i++ {
		idx := i
		g.Go(func() error {
			if idx%2 == 0 {
				panic("boom")
			}
			done.Add(1)
			return nil
		})
	}

	_, timeout := waitWithTimeout(g, 3*time.Second)
	if timeout {
		t.Fatal("Wait() 超时，worker 被 panic 破坏")
	}
	// 第一个任务 panic 会 cancel ctx，后续任务可能被跳过，这里只要求不死锁不崩溃
	t.Logf("正常完成的任务数: %d", done.Load())
}

// TestGoAfterWaitDoesNotBlock 回归：Wait() 之后再 Go() 不能永久阻塞
func TestGoAfterWaitDoesNotBlock(t *testing.T) {
	g := New(4)
	g.Go(func() error { return nil })

	if err := g.Wait(); err != nil {
		t.Fatalf("首次 Wait() 不应有 error: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		g.Go(func() error { return nil })
		g.Wait()
	}()

	select {
	case <-done:
		// 正确：不阻塞
	case <-time.After(3 * time.Second):
		t.Fatal("Wait() 之后调用 Go() 造成永久阻塞（goroutine 泄漏）")
	}
}

// TestWaitIdempotent 多次 Wait 应该返回一致的结果且不 panic
func TestWaitIdempotent(t *testing.T) {
	expected := errors.New("task failed")

	g := New(4)
	g.Go(func() error { return expected })

	first := g.Wait()
	second := g.Wait()
	third := g.Wait()

	if !errors.Is(first, expected) {
		t.Errorf("第 1 次 Wait 期望 %v，实际 %v", expected, first)
	}
	if first != second || second != third {
		t.Errorf("多次 Wait 返回值不一致: %v / %v / %v", first, second, third)
	}
}

// TestWaitCancelsContext 回归：Wait 返回后必须 cancel ctx，避免 context 泄漏
func TestWaitCancelsContext(t *testing.T) {
	g, ctx := WithContext(context.Background(), 4)
	g.Go(func() error { return nil })

	if err := g.Wait(); err != nil {
		t.Fatalf("不应有 error: %v", err)
	}

	select {
	case <-ctx.Done():
		// 正确
	case <-time.After(time.Second):
		t.Fatal("Wait 返回后 ctx 应该被 cancel")
	}
}

// TestInvalidPoolSize 池大小非法时应回落到默认值而不是崩溃
func TestInvalidPoolSize(t *testing.T) {
	for _, size := range []int{0, -1, -100} {
		g, _ := WithContext(context.Background(), size)
		if g.pool == nil {
			t.Fatalf("poolSize=%d 时池不应为 nil", size)
		}

		var n atomic.Int32
		for i := 0; i < 10; i++ {
			g.Go(func() error { n.Add(1); return nil })
		}

		if _, timeout := waitWithTimeout(g, 3*time.Second); timeout {
			t.Fatalf("poolSize=%d 时 Wait 超时", size)
		}
		if n.Load() != 10 {
			t.Errorf("poolSize=%d 期望执行 10 个任务，实际 %d", size, n.Load())
		}
	}
}

// TestNilTaskIgnored 提交 nil 任务不应导致 Wait 阻塞
func TestNilTaskIgnored(t *testing.T) {
	g := New(4)
	g.Go(nil)
	g.Go(func() error { return nil })

	if _, timeout := waitWithTimeout(g, 3*time.Second); timeout {
		t.Fatal("提交 nil 任务后 Wait 阻塞")
	}
}

// TestNoGoroutineLeakOnUpstreamCancel 上游 ctx 取消后，Wait 仍应能正常返回
func TestNoGoroutineLeakOnUpstreamCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	g, _ := WithContext(ctx, 8)

	for i := 0; i < 50; i++ {
		g.Go(func() error {
			time.Sleep(20 * time.Millisecond)
			return nil
		})
	}
	cancel()

	if _, timeout := waitWithTimeout(g, 5*time.Second); timeout {
		t.Fatal("上游 cancel 后 Wait 未能返回")
	}
}

