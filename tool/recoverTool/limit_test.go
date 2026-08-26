package recoverTool

import (
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

type countingEmailSvc struct {
	mx       sync.Mutex
	subjects []string
}

func (c *countingEmailSvc) Send(toL []string, subject, content string) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.subjects = append(c.subjects, subject)
	return nil
}

func (c *countingEmailSvc) count() int {
	c.mx.Lock()
	defer c.mx.Unlock()
	return len(c.subjects)
}

func newTestLimit(cooldown time.Duration) *Limit {
	return &Limit{pool: make(map[string]*limitEntry), cooldown: cooldown}
}

func newTestRecoverTool(svc EmailSvc, cooldown time.Duration) *RecoverTool {
	return &RecoverTool{
		service:       "test",
		env:           "unit",
		emailReceiver: []string{"a@example.com"},
		emailSvc:      svc,
		limit:         newTestLimit(cooldown),
		lg:            logrus.New(),
	}
}

// HandleRecover 可以被 gin/grpc 拦截器直接调用，此时 emailSvc 可能为 nil。
// 由于调用方通常放在裸 goroutine 里，这里 panic 会直接崩掉整个进程。
func TestHandleRecoverWithNilEmailSvc(t *testing.T) {
	r := &RecoverTool{
		service: "test",
		env:     "unit",
		limit:   newTestLimit(time.Hour),
		lg:      logrus.New(),
	}

	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("HandleRecover 不应该 panic, got: %v", p)
		}
	}()

	r.HandleRecover("boom", "stack")
}

func TestHandleRecoverWithEmptyReceiver(t *testing.T) {
	svc := &countingEmailSvc{}
	r := newTestRecoverTool(svc, time.Hour)
	r.emailReceiver = nil

	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("HandleRecover 不应该 panic, got: %v", p)
		}
	}()

	r.HandleRecover("boom", "stack")

	if svc.count() != 0 {
		t.Errorf("收件人为空时不应发送邮件, sent=%d", svc.count())
	}
}

func TestHandleRecoverSuppressesWithinCooldown(t *testing.T) {
	svc := &countingEmailSvc{}
	r := newTestRecoverTool(svc, time.Hour)

	for i := 0; i < 10; i++ {
		r.HandleRecover("boom", "stack")
	}

	if svc.count() != 1 {
		t.Errorf("冷却期内只应发送 1 封邮件, sent=%d", svc.count())
	}
}

func TestShouldSendSuppressedCountExcludesSentOne(t *testing.T) {
	l := newTestLimit(50 * time.Millisecond)

	if ok, sup := l.ShouldSend("boom"); !ok || sup != 0 {
		t.Fatalf("首次应发送且无抑制, got (%v, %d)", ok, sup)
	}
	// 冷却期内的 3 次都应被抑制
	for i := 0; i < 3; i++ {
		if ok, _ := l.ShouldSend("boom"); ok {
			t.Fatalf("第 %d 次应被抑制", i)
		}
	}

	time.Sleep(60 * time.Millisecond)

	ok, sup := l.ShouldSend("boom")
	if !ok {
		t.Fatal("冷却结束后应发送")
	}
	if sup != 3 {
		t.Errorf("抑制次数应为 3（不含本次发送）, got %d", sup)
	}

	// 发送后计数归零
	if ok, _ := l.ShouldSend("boom"); ok {
		t.Fatal("刚发送完应被抑制")
	}
	if got := l.pool["boom"].suppressed; got != 1 {
		t.Errorf("suppressed 应重新从 1 开始, got %d", got)
	}
}

func TestShouldSendIsAtomic(t *testing.T) {
	l := newTestLimit(time.Hour)

	var wg sync.WaitGroup
	var mx sync.Mutex
	sent := 0

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ok, _ := l.ShouldSend("boom"); ok {
				mx.Lock()
				sent++
				mx.Unlock()
			}
		}()
	}
	wg.Wait()

	if sent != 1 {
		t.Errorf("并发下只应发送 1 次, sent=%d", sent)
	}
}

func TestShouldSendDistinctKeysAreIndependent(t *testing.T) {
	l := newTestLimit(time.Hour)

	if ok, _ := l.ShouldSend("a"); !ok {
		t.Error("key a 首次应发送")
	}
	if ok, _ := l.ShouldSend("b"); !ok {
		t.Error("key b 首次应发送")
	}
	if ok, _ := l.ShouldSend("a"); ok {
		t.Error("key a 第二次应被抑制")
	}
}
