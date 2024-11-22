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

		secondK: fmt.Sprintf("%s:%s:counter:second", opt.Service, opt.Prefix),
		minuteK: fmt.Sprintf("%s:%s:counter:minute", opt.Service, opt.Prefix),
		hourK:   fmt.Sprintf("%s:%s:counter:hour", opt.Service, opt.Prefix),
		dayK:    fmt.Sprintf("%s:%s:counter:day", opt.Service, opt.Prefix),
		weekK:   fmt.Sprintf("%s:%s:counter:week", opt.Service, opt.Prefix),
		monthK:  fmt.Sprintf("%s:%s:counter:month", opt.Service, opt.Prefix),
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
		switch v {
		case c.secondK:
			k := fmt.Sprintf("%s:%s", v, tool.Second())
			c.redis.ZIncrBy(ctx, k, float64(score), itemS)
			c.redis.Expire(ctx, k, c.ttl)
		case c.minuteK:
			k := fmt.Sprintf("%s:%s", v, tool.Minute())
			c.redis.ZIncrBy(ctx, k, float64(score), itemS)
			c.redis.Expire(ctx, k, c.ttl)
		case c.hourK:
			k := fmt.Sprintf("%s:%s", v, tool.Hour())
			c.redis.ZIncrBy(ctx, k, float64(score), itemS)
			c.redis.Expire(ctx, k, c.ttl)
		case c.dayK:
			k := fmt.Sprintf("%s:%s", v, tool.Day())
			c.redis.ZIncrBy(ctx, k, float64(score), itemS)
			c.redis.Expire(ctx, k, c.ttl)
		case c.weekK:
			k := fmt.Sprintf("%s:%s", v, tool.Week())
			c.redis.ZIncrBy(ctx, k, float64(score), itemS)
			c.redis.Expire(ctx, k, c.ttl)
		case c.monthK:
			k := fmt.Sprintf("%s:%s", v, tool.Month())
			c.redis.ZIncrBy(ctx, k, float64(score), itemS)
			c.redis.Expire(ctx, k, c.ttl)
		}
	}
}
