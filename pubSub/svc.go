package pubSub

import (
	"context"
	"errors"
	"fmt"
	"os"

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

//var rc redis.UniversalClient
//var service string

type Svc struct {
	rc      redis.UniversalClient
	service string
}

func NewSvc(optChain ...*Option) *Svc {
	if len(optChain) == 0 {
		panic("opt is nil")
	}

	opt := optChain[0]

	if opt == nil {
		panic("opt is nil")
	}

	if opt.RedisClient == nil {
		panic("opt redis client is nil")
	}

	if opt.Service == "" {
		panic("opt service is empty")
	}

	//if svc != nil {
	//	panic("redis client already initialized")
	//}

	svc := &Svc{
		rc:      opt.RedisClient,
		service: opt.Service,
	}

	return svc
}

func (s *Svc) Publish(ctx context.Context, topic string, msg interface{}) error {
	switch msg.(type) {
	case string:
		return s.rc.Publish(ctx, topic, msg.(string)).Err()
	default:
		b, err := tool.Marshal(msg)
		if err != nil {
			return err
		}
		return s.rc.Publish(ctx, topic, b).Err()
	}
}

func (s *Svc) RunConsumer(topic string, f func(msg string) error) error {

	if s.rc == nil {
		return fmt.Errorf("redis client is nil, please call Init first")
	}

	c := tool.NewCtxBK()

	sub := s.rc.Subscribe(c, topic)
	defer sub.Close()

	if _, err := sub.Receive(c); err != nil {
		panic(err)
	}

	msgCh := sub.Channel()

	for {
		ctx := tool.NewCtxBK()
		lg := log.CtxLogger(ctx).WithFields(logrus.Fields{"topic": topic, "service": s.service})

		m := <-msgCh
		if m == nil {
			lg.Errorf("msg is nil")
			continue
		}

		d := m.Payload

		lockK := fmt.Sprintf("%s:pub-sub:consumer:lock:%s:%s", s.service, topic, tool.MD5(d))
		_, err := lock.GetLocker(s.rc).Obtain(ctx, lockK, 10)
		if err != nil {
			if errors.Is(err, redislock.ErrNotObtained) {
				if os.Getenv("DEBUD") == "TRUE" {
					lg.Infof("lock %s not obtained, msg: %s", lockK, d)
				}
				continue
			}
		}

		if os.Getenv("DEBUD") == "TRUE" {
			lg.Infof("recv msg: %s", d)
		}

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
