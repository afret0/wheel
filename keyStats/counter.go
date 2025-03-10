package keyStats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afret0/wheel/tool"
	"github.com/redis/go-redis/v9"
	"time"
)

type Option struct {
	Service string
	Prefix  string
	TTL     time.Duration

	Redis redis.UniversalClient
}

type Counter struct {
	redis redis.UniversalClient
	ttl   time.Duration

	secondK string
	minuteK string
	hourK   string
	dayK    string
	weekK   string
	monthK  string
}

type Item struct {
	Name  string `json:"name"`
	Extra string `json:"extra"`
}

func (i *Item) Marshal() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func NewCounter(opt *Option) *Counter {
	if opt.Service == "" {
		panic("service can not be empty")
	}
	if opt.TTL <= 0 {
		opt.TTL = time.Duration(24*3) * time.Hour
	}
	if opt.Prefix == "" {
		opt.Prefix = opt.Service
	}

	return &Counter{
		redis: opt.Redis,
		ttl:   opt.TTL,

		secondK: fmt.Sprintf("%s:%s:keyStats:second", opt.Service, opt.Prefix),
		minuteK: fmt.Sprintf("%s:%s:keyStats:minute", opt.Service, opt.Prefix),
		hourK:   fmt.Sprintf("%s:%s:keyStats:hour", opt.Service, opt.Prefix),
		dayK:    fmt.Sprintf("%s:%s:keyStats:day", opt.Service, opt.Prefix),
		weekK:   fmt.Sprintf("%s:%s:keyStats:week", opt.Service, opt.Prefix),
		monthK:  fmt.Sprintf("%s:%s:keyStats:month", opt.Service, opt.Prefix),
	}
}

func (c *Counter) Incr(ctx context.Context, item *Item) {
	go func() {
		c.incrBy(ctx, item, 1)
	}()
}

func (c *Counter) IncrBy(ctx context.Context, item *Item, score int64) {
	go func() {
		c.incrBy(ctx, item, score)
	}()
}

func (c *Counter) incrBy(ctx context.Context, item *Item, score int64) {
	itemS := item.Marshal()
	for _, v := range []string{c.secondK, c.minuteK, c.hourK, c.dayK, c.weekK, c.monthK} {
		k := ""
		switch v {
		case c.secondK:
			k = fmt.Sprintf("%s:%s", v, tool.Second())
		case c.minuteK:
			k = fmt.Sprintf("%s:%s", v, tool.Minute())
		case c.hourK:
			k = fmt.Sprintf("%s:%s", v, tool.Hour())
		case c.dayK:
			k = fmt.Sprintf("%s:%s", v, tool.Day())
		case c.weekK:
			k = fmt.Sprintf("%s:%s", v, tool.Week())
		case c.monthK:
			k = fmt.Sprintf("%s:%s", v, tool.Month())
		}

		kt := fmt.Sprintf("%s:total", k)
		c.redis.ZIncrBy(ctx, k, float64(score), itemS)
		c.redis.IncrBy(ctx, kt, score)
		c.redis.Expire(ctx, k, c.ttl)
		c.redis.Expire(ctx, kt, c.ttl)
	}
}
