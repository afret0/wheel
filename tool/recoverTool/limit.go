package recoverTool

import (
	"sync"
	"time"
)

type limitEntry struct {
	count      int64
	lastSent   time.Time
	suppressed int64
}

type Limit struct {
	mx       sync.Mutex
	pool     map[string]*limitEntry
	cooldown time.Duration
}

var limit *Limit

const defaultCooldown = 30 * time.Minute

func GetLimit() *Limit {
	if limit != nil {
		return limit
	}
	limit = &Limit{
		mx:       sync.Mutex{},
		pool:     make(map[string]*limitEntry),
		cooldown: defaultCooldown,
	}

	go limit.clearPoolLoop()

	return limit
}

// clearPoolLoop 定期清理长时间未出现的 key，避免内存泄漏
func (l *Limit) clearPoolLoop() {
	for range time.Tick(l.cooldown) {
		l.mx.Lock()
		for k, v := range l.pool {
			if time.Since(v.lastSent) > l.cooldown*2 {
				delete(l.pool, k)
			}
		}
		l.mx.Unlock()
	}
}

// ShouldSend 原子地检查是否应该发送邮件
// 返回: 是否发送, 自上次发送以来被抑制的次数
func (l *Limit) ShouldSend(key string) (bool, int64) {
	l.mx.Lock()
	defer l.mx.Unlock()

	entry, exists := l.pool[key]
	if !exists {
		l.pool[key] = &limitEntry{
			count:      1,
			lastSent:   time.Now(),
			suppressed: 0,
		}
		return true, 0
	}

	entry.count++
	entry.suppressed++

	if time.Since(entry.lastSent) >= l.cooldown {
		suppressed := entry.suppressed
		entry.lastSent = time.Now()
		entry.suppressed = 0
		return true, suppressed
	}

	return false, 0
}
