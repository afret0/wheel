package pubSub

import (
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/afret0/wheel/tool"
)

func Test_pubSub(t *testing.T) {
	RC := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"r-bp1kvud328x48r9xp6pd.redis.rds.aliyuncs.com:6379"},
		Username: "kiwi0621",
		Password: "Qiyiguo0621",
	})

	pb := NewSvc(&Option{RedisClient: RC, Service: "test"})

	topic := "test-topic"

	go func() {
		f := func(msg string) error {
			fmt.Printf("%s\n", msg)
			return nil
		}

		err := pb.RunConsumer(topic, f)
		if err != nil {
			t.Errorf("%v", err)
		}
	}()

	for {
		ctx := tool.NewCtxBK()
		err := pb.Publish(ctx, topic, fmt.Sprintf("%s", time.Now().String()))
		if err != nil {
			t.Errorf("publish error: %v", err)
		}
		time.Sleep(time.Second)
	}
}
