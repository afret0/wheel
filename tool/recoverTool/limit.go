package recoverTool

import (
	"sync"
	"time"
)

type Limit struct {
	mx   sync.Mutex
	pool map[string]int64
}

var limit *Limit

func GetLimit() *Limit {
	if limit != nil {
		return limit
	}
	limit = &Limit{
		mx:   sync.Mutex{},
		pool: make(map[string]int64),
	}

	go limit.clearPoolLoop()

	return limit
}

func (l *Limit) clearPoolLoop() {
	for _ = range time.Tick(time.Minute * 5) {
		l.mx.Lock()
		l.pool = make(map[string]int64)
		l.mx.Unlock()
	}
}

func (l *Limit) Incr(key string) {
	l.mx.Lock()
	defer l.mx.Unlock()
	if _, ok := l.pool[key]; !ok {
		l.pool[key] = 0
	}
	l.pool[key]++
}

func (l *Limit) Count(key string) int64 {
	l.mx.Lock()
	defer l.mx.Unlock()
	if _, ok := l.pool[key]; !ok {
		return 0
	}
	return l.pool[key]
}
