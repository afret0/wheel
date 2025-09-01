package pubSub

import (
	"context"
	"errors"
	"fmt"

	"github.com/afret0/wheel/lock"
	"github.com/afret0/wheel/log"
	"github.com/afret0/wheel/tool"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Option struct {
	RedisClient redis.UniversalClient
	Service     string
}

var rc redis.UniversalClient
var service string

func Init(opt *Option) {
	if opt == nil {
		panic("opt is nil")
	}

	if opt.RedisClient == nil {
		panic("opt redis client is nil")
	}

	if opt.Service == "" {
		panic("opt service is empty")
	}

	if rc != nil {
		panic("redis client already initialized")
	}

	rc = opt.RedisClient
	service = opt.Service
}

func Publish(ctx context.Context, topic string, msg interface{}) error {
	switch msg.(type) {
	case string:
		return rc.Publish(ctx, topic, msg.(string)).Err()
	default:
		s, err := tool.Marshal(msg)
		if err != nil {
			return err
		}
		return rc.Publish(ctx, topic, s).Err()
	}
}

func RunConsumer(topic string, f func(msg string) error) {

	c := tool.NewCtxBK()

	sub := rc.Subscribe(c, topic)
	defer sub.Close()

	if _, err := sub.Receive(c); err != nil {
		panic(err)
	}

	msgCh := sub.Channel()

	for {
		ctx := tool.NewCtxBK()
		lg := log.CtxLogger(ctx).WithFields(logrus.Fields{"topic": topic, "service": service})

		m := <-msgCh
		if m == nil {
			lg.Errorf("msg is nil")
			continue
		}

		d := m.Payload

		lockK := fmt.Sprintf("%s:pub-sub:consumer:lock:%s:%s", service, topic, tool.MD5(d))
		_, err := lock.GetLocker(rc).Obtain(ctx, lockK, 10)
		if err != nil {
			if errors.Is(err, redislock.ErrNotObtained) {

			}
		}

		lg.Infof("recv msg: %s", d)

		err = f(d)
		if err != nil {
			lg.Errorf("err: %s", err)
		}

		//select {
		//case m := <-msgCh:
		//	ctx := tool.NewCtxBK()
		//	lg := log.CtxLogger(ctx).WithFields(logrus.Fields{})
		//
		//	if m == nil {
		//		lg.Errorf("msg is nil")
		//		continue
		//	}
		//}
	}
}
