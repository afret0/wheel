package timeMock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/afret0/wheel/log"
	"github.com/afret0/wheel/tool"
	"github.com/afret0/wheel/tool/timeTool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var mx sync.Mutex
var rc redis.UniversalClient
var k string

//var optTag bool

const debugTag = "TIME_TOOL_DEBUG"

type Option struct {
	RedisClient redis.UniversalClient
	KeyPrefix   string
}

func SetOption(opt *Option) error {
	mx.Lock()
	defer mx.Unlock()

	if opt.RedisClient == nil {
		return fmt.Errorf("redis client is nil")
	}

	if opt.KeyPrefix == "" {
		return fmt.Errorf("key prefix is empty")
	}

	rc = opt.RedisClient
	k = fmt.Sprintf("%s:time_tool:ts", opt.KeyPrefix)

	return nil
}

func SetTime(ctx context.Context, ts int64) error {

	if !tool.Debug(debugTag) {
		return fmt.Errorf("debug mode not enabled")
	}

	mx.Lock()
	defer mx.Unlock()

	if ts == 0 {
		return fmt.Errorf("timestamp is zero")
	}

	err := rc.Set(ctx, k, ts, time.Hour*24*3).Err()
	return err
}

func now() time.Time {
	return time.Now().In(timeTool.Location())
}

func Now(ctx context.Context) time.Time {
	lg := log.CtxLogger(ctx).WithFields(logrus.Fields{})

	if !tool.Debug(debugTag) {
		return now()
	}

	if rc == nil {
		panic("tm redis client is nil")
	}

	ts, err := rc.Get(ctx, k).Int64()
	if err != nil {
		lg.Errorf("err: %s", err)
		return now()
	}

	if ts == 0 {
		lg.Errorf("timestamp is zero")
		return now()
	}

	return time.UnixMilli(ts).In(timeTool.Location())
}
