package timeMock

import (
	"os"
	"testing"
	"time"

	"github.com/afret0/wheel/tool"
	"github.com/redis/go-redis/v9"
)

func Test_time(t *testing.T) {
	RC := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"r-bp1kvud328x48r9xp6pd.redis.rds.aliyuncs.com:6379"},
		Username: "kiwi0621",
		Password: "Qiyiguo0621",
	})

	ctx := tool.NewCtxBK()
	now := time.Now()

	t.Logf("now: %s", now.String())
	t.Logf("now: %d", now.Day())

	os.Setenv("TIME_TOOL_DEBUG", "true")

	err := SetOption(&Option{RedisClient: RC, KeyPrefix: "test"})
	if err != nil {
		t.Errorf("SetOption error: %v", err)
	}

	err = SetTime(ctx, 1735689600000)
	if err != nil {
		t.Errorf("SetTime error: %v", err)
	}

	now1 := Now(ctx)

	t.Logf("now: %s, ts: %d", now1.String(), now1.UnixMilli())

}
