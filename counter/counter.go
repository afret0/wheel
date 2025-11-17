package counter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Counter struct {
	redis  redis.UniversalClient
	prefix string
	ttl    time.Duration
}

type Option struct {
	RedisClient redis.UniversalClient
	Prefix      string
	TTL         time.Duration
}

type Opt = Option

var counter *Counter
var opt *Option

func SetOption(o *Option) {
	opt = o
}

func GetCounter() *Counter {
	if counter != nil {
		return counter
	}

	counter := NewCounter(opt)

	return counter
}

func NewCounter(opt *Option) *Counter {
	if opt == nil {
		panic("option is nil")
	}
	if opt.RedisClient == nil {
		panic("redis client is nil")
	}
	if opt.Prefix == "" {
		panic("prefix is empty")
	}
	if opt.TTL == 0 {
		opt.TTL = time.Hour * 24 * 100
	}

	return &Counter{
		redis:  opt.RedisClient,
		prefix: opt.Prefix,
		ttl:    opt.TTL,
	}
}

func (c *Counter) Key(k string) string {
	return fmt.Sprintf("%s:counter:stats:%s", c.prefix, k)
}

func (c *Counter) Incr(ctx context.Context, key string, numChain ...int64) (int64, error) {
	//lg := log.CtxLogger(ctx).WithFields(logrus.Fields{})

	num := int64(1)
	if len(numChain) > 0 {
		num = numChain[0]
	}

	p := c.redis.Pipeline()

	r, err := p.IncrBy(ctx, c.Key(key), num).Result()
	if err != nil {
		//lg.Errorf("err: %s", err)
		return 0, err
	}

	_, err = p.Expire(ctx, c.Key(key), c.ttl).Result()
	if err != nil {
		//lg.Errorf("err: %s", err)
		return 0, err
	}

	_, err = p.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return r, nil
}

func (c *Counter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	_, err := c.redis.Expire(ctx, c.Key(key), ttl).Result()
	return err
}

func (c *Counter) Get(ctx context.Context, key string) (int64, error) {
	r, err := c.redis.Get(ctx, c.Key(key)).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}

	return r, nil
}

func (c *Counter) IsExceeded(ctx context.Context, key string, limit int64) (bool, error) {
	count, err := c.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if count >= limit {
		return true, nil
	}
	return false, nil
}
