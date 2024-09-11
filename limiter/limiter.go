package limiter

import (
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

func NewLimiter(rdb redis.UniversalClient) *redis_rate.Limiter {
	limiter := redis_rate.NewLimiter(rdb)
	return limiter
}

func PerDuration(rate int, duration time.Duration) redis_rate.Limit {
	return redis_rate.Limit{
		Rate:   rate,
		Burst:  rate,
		Period: duration / time.Duration(rate),
	}
}
