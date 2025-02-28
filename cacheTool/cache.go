package cacheTool

import (
	"context"
	"errors"
	"github.com/afret0/wheel/log"
	redisCache "github.com/go-redis/cache/v9"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Option struct {
	NoCacheMongoNoDocuments bool
	SetXX                   bool
	SetNX                   bool
	SkipLocalCache          bool
}

type cacheResult[T any] struct {
	Data T      `json:"data"`
	Err  string `json:"err"`
}

func WithCache[T any](
	ctx context.Context,
	cache *redisCache.Cache,
	cacheKey string,
	ttl time.Duration,
	fetchFunc func(ctx context.Context) (T, error),
	optChains ...*Option,
) (T, error) {
	lg := log.CtxLogger(ctx).WithFields(logrus.Fields{"cacheKey": cacheKey})
	result := cacheResult[T]{}

	opt := Option{}
	if len(optChains) > 0 && optChains[0] != nil {
		opt = *optChains[0]
	}

	err := cache.Once(&redisCache.Item{
		Ctx:   ctx,
		Key:   cacheKey,
		Value: &result,
		TTL:   ttl,
		Do: func(item *redisCache.Item) (interface{}, error) {
			data, err := fetchFunc(ctx)
			if err != nil {
				if opt.NoCacheMongoNoDocuments && errors.Is(err, mongo.ErrNoDocuments) {
					return &cacheResult[T]{
						Err: err.Error(),
					}, nil
				}
				lg.Errorf("err: %s", err)
				return nil, err
			}

			return &cacheResult[T]{
				Data: data,
			}, nil
		},
		SetNX:          opt.SetNX,
		SetXX:          opt.SetXX,
		SkipLocalCache: opt.SkipLocalCache,
	})

	if err != nil {
		lg.Errorf("Cache error: %s", err)
		return result.Data, err
	}

	// If we cached an error, return it
	if result.Err != "" {
		return result.Data, errors.New(result.Err)
	}

	return result.Data, nil
}
