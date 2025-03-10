package keyStats

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func Test_counter(t *testing.T) {
	RC := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"r-bp1kvud328x48r9xp6pd.redis.rds.aliyuncs.com:6379"},
		Username: "kiwi0621",
		Password: "Qiyiguo0621",
	})

	C := NewCounter(&Option{
		Service: "wheel:test",
		Prefix:  "counter",
		TTL:     time.Duration(10) * time.Minute,
		Redis:   RC,
	})

	for i := 0; i < 100; i++ {
		C.Incr(context.Background(), &Item{Name: "test"})
		C.IncrBy(context.Background(), &Item{Name: "test2"}, 10)
	}

	time.Sleep(10 * time.Second)
}
