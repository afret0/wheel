package antsGroup

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// TestNew 测试 New 函数
func TestNew(t *testing.T) {
	g := New()
	if g == nil {
		t.Fatal("New() returned nil")
	}
	if g.pool == nil {
		t.Error("pool is nil")
	}
	if g.ctx == nil {
		t.Error("ctx is nil")
	}
	if g.cancel == nil {
		t.Error("cancel is nil")
	}
}

// TestNew_WithPoolSize 测试带自定义池大小的 New 函数
func TestNew_WithPoolSize(t *testing.T) {
	poolSize := 10
	g := New(poolSize)
	if g == nil {
		t.Fatal("New() returned nil")
	}

	// 测试池是否可用
	var counter int32
	for i := 0; i < 20; i++ {
		g.Go(func() error {
			atomic.AddInt32(&counter, 1)
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		t.Errorf("Wait() returned error: %v", err)
	}

	if atomic.LoadInt32(&counter) != 20 {
		t.Errorf("expected 20 tasks executed, got %d", counter)
	}
}

// TestWithContext 测试 WithContext 函数
func TestWithContext(t *testing.T) {
	ctx := context.Background()
	poolSize := 5
	g, derivedCtx := WithContext(ctx, poolSize)

	if g == nil {
		t.Fatal("WithContext() returned nil group")
	}
	if derivedCtx == nil {
		t.Fatal("WithContext() returned nil context")
	}
	if g.pool == nil {
		t.Error("pool is nil")
	}
	if g.ctx == nil {
		t.Error("ctx is nil")
	}
}

// TestGroup_Go_Success 测试成功执行任务
func TestGroup_Go_Success(t *testing.T) {
	g := New()

	var counter int32
	taskCount := 100

	for i := 0; i < taskCount; i++ {
		g.Go(func() error {
			atomic.AddInt32(&counter, 1)
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		t.Errorf("Wait() returned unexpected error: %v", err)
	}

	if atomic.LoadInt32(&counter) != int32(taskCount) {
		t.Errorf("expected %d tasks executed, got %d", taskCount, counter)
	}
}

// TestGroup_Go_WithError 测试任务返回错误
func TestGroup_Go_WithError(t *testing.T) {
	g := New()

	expectedErr := errors.New("task error")
	var counter int32

	// 提交多个任务，其中一个会失败
	for i := 0; i < 10; i++ {
		idx := i
		g.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt32(&counter, 1)
			if idx == 5 {
				return expectedErr
			}
			return nil
		})
	}

	err := g.Wait()
	if err == nil {
		t.Error("Wait() should return error but got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

// TestGroup_Go_MultipleErrors 测试多个任务返回错误（应该返回第一个错误）
func TestGroup_Go_MultipleErrors(t *testing.T) {
	g := New(10)

	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	var errorReceived error

	g.Go(func() error {
		time.Sleep(10 * time.Millisecond)
		return err1
	})

	g.Go(func() error {
		time.Sleep(20 * time.Millisecond)
		return err2
	})

	errorReceived = g.Wait()

	if errorReceived == nil {
		t.Fatal("Wait() should return error but got nil")
	}

	// 应该返回第一个错误
	if errorReceived != err1 && errorReceived != err2 {
		t.Errorf("expected err1 or err2, got %v", errorReceived)
	}
}

// TestGroup_ContextCancellation 测试 context 取消功能
func TestGroup_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	g, _ := WithContext(ctx, 10)

	var startedTasks int32
	var completedTasks int32

	// 提交多个长时间运行的任务
	for i := 0; i < 20; i++ {
		g.Go(func() error {
			atomic.AddInt32(&startedTasks, 1)
			time.Sleep(100 * time.Millisecond)
			atomic.AddInt32(&completedTasks, 1)
			return nil
		})
	}

	// 短时间后取消 context
	time.Sleep(50 * time.Millisecond)
	cancel()

	err := g.Wait()
	if err != nil {
		t.Logf("Wait() returned error: %v", err)
	}

	// 由于 context 被取消，部分任务可能不会完成
	started := atomic.LoadInt32(&startedTasks)
	completed := atomic.LoadInt32(&completedTasks)
	t.Logf("Started: %d, Completed: %d", started, completed)
}

// TestGroup_ErrorTriggersContextCancel 测试错误触发 context 取消
func TestGroup_ErrorTriggersContextCancel(t *testing.T) {
	g, ctx := WithContext(context.Background(), 10)

	expectedErr := errors.New("trigger error")
	var tasksAfterError int32

	// 第一个任务立即失败
	g.Go(func() error {
		return expectedErr
	})

	// 后续任务应该因为 context 取消而跳过
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			time.Sleep(50 * time.Millisecond)
			// 检查 context 是否被取消
			select {
			case <-ctx.Done():
				return nil
			default:
				atomic.AddInt32(&tasksAfterError, 1)
				return nil
			}
		})
	}

	err := g.Wait()
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	// context 应该被取消
	select {
	case <-ctx.Done():
		// 正确：context 被取消了
	default:
		t.Error("context should be cancelled after error")
	}
}

// TestGroup_Concurrent 测试并发安全性
func TestGroup_Concurrent(t *testing.T) {
	g := New(20)

	var counter int32
	taskCount := 1000

	for i := 0; i < taskCount; i++ {
		g.Go(func() error {
			atomic.AddInt32(&counter, 1)
			time.Sleep(1 * time.Millisecond)
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		t.Errorf("Wait() returned error: %v", err)
	}

	if atomic.LoadInt32(&counter) != int32(taskCount) {
		t.Errorf("expected %d tasks executed, got %d", taskCount, counter)
	}
}

// TestGroup_EmptyGroup 测试空的 Group（没有提交任何任务）
func TestGroup_EmptyGroup(t *testing.T) {
	g := New()

	err := g.Wait()
	if err != nil {
		t.Errorf("Wait() on empty group should return nil, got %v", err)
	}
}

// TestGroup_Wait_MultipleCalls 测试多次调用 Wait（虽然不推荐）
func TestGroup_Wait_MultipleCalls(t *testing.T) {
	g := New()

	g.Go(func() error {
		return nil
	})

	err1 := g.Wait()
	if err1 != nil {
		t.Errorf("First Wait() returned error: %v", err1)
	}

	// 注意：第二次调用 Wait 可能会因为 pool 已 Release 而有不同行为
	// 这里只是测试不会 panic
	err2 := g.Wait()
	if err2 != nil {
		t.Logf("Second Wait() returned error: %v", err2)
	}
}

// BenchmarkGroup_Go 基准测试
func BenchmarkGroup_Go(b *testing.B) {
	g := New(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Go(func() error {
			return nil
		})
	}

	_ = g.Wait()
}

// BenchmarkGroup_GoWithWork 带实际工作的基准测试
func BenchmarkGroup_GoWithWork(b *testing.B) {
	g := New(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Go(func() error {
			// 模拟一些工作
			var sum int
			for j := 0; j < 100; j++ {
				sum += j
			}
			return nil
		})
	}

	_ = g.Wait()
}
