package lock

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type Locker struct {
	Locker *redislock.Client
}

var locker *Locker

func GetLocker(rdb redis.UniversalClient) *Locker {
	if locker != nil {
		return locker
	}
	locker = new(Locker)
	locker.Locker = redislock.New(rdb)

	return locker
}

func (l *Locker) Obtain(ctx context.Context, key string, ttl int, opt ...*redislock.Options) (*redislock.Lock, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	t := time.Second * time.Duration(5)
	if ttl != 0 {
		t = time.Duration(ttl) * time.Second
	}

	optN := new(redislock.Options)
	if len(opt) != 0 {
		optN = opt[0]
	}

	return l.Locker.Obtain(ctx, key, t, optN)
}

func (l *Locker) ObtainWaitRetry(ctx context.Context, key string, ttl int, retryCount int) (*redislock.Lock, error) {
	backoff := redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), retryCount)
	return l.Locker.Obtain(ctx, key, time.Duration(ttl)*time.Second, &redislock.Options{
		RetryStrategy: backoff,
	})
	//timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(retryDelay*retryCount)*time.Second+1)
	//defer cancel()
	//
	//for {
	//	select {
	//	case <-timeoutCtx.Done():
	//		// If the context is done (because of the timeout), return an error
	//		return nil, errors.New("timeout while trying to obtain the lock")
	//	default:
	//		// Try to obtain the lock
	//		lock, err := l.Locker.Obtain(ctx, key, time.Duration(ttl)*time.Second, nil)
	//		if err != nil {
	//			if errors.Is(err, redislock.ErrNotObtained) {
	//				// If the lock is not obtained, wait for a while before trying again
	//				time.Sleep(time.Duration(retryDelay) * time.Second)
	//				continue
	//			}
	//			// If there is another error, return it
	//			return nil, err
	//		}
	//
	//		// If the lock is obtained, return it
	//		return lock, nil
	//	}
	//}
}

func (l *Locker) ObtainWaitExponentialRetry(ctx context.Context, key string, ttl int, maxWaitTime int) (*redislock.Lock, error) {
	backoff := redislock.LimitRetry(redislock.ExponentialBackoff(100*time.Millisecond, time.Duration(maxWaitTime)*time.Millisecond), -1)
	return l.Locker.Obtain(ctx, key, time.Duration(ttl)*time.Second, &redislock.Options{
		RetryStrategy: backoff,
	})
}
