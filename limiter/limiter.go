package limiter

import (
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-redis/redis_rate/v10"
)

// deprecated
func NewLimiter(rdb redis.UniversalClient) *redis_rate.Limiter {
	limiter := redis_rate.NewLimiter(rdb)
	return limiter
}

// deprecated
func PerDuration(rate int, duration time.Duration) redis_rate.Limit {
	return redis_rate.Limit{
		Rate:   rate,
		Burst:  rate,
		Period: duration,
	}
}
