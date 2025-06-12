package cacheTool

import (
	"context"
	"errors"
	"time"

	redisCache "github.com/go-redis/cache/v9"
	"go.mongodb.org/mongo-driver/mongo"
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
	//ErrType string `json:"errType"`
}

//var errorMapping = map[string]error{
//	"mongo.ErrNoDocuments": mongo.ErrNoDocuments,
//	// 添加其他需要识别的错误...
//}

type SliceWrapper[T any] struct {
	L []*T `json:"l" msgpack:"l"`
}

func WithCache[T any](
	ctx context.Context,
	cache *redisCache.Cache,
	cacheKey string,
	ttl time.Duration,
	fetchFunc func(ctx context.Context) (T, error),
	optChains ...*Option,
) (T, error) {
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
				if !opt.NoCacheMongoNoDocuments && errors.Is(err, mongo.ErrNoDocuments) {
					return &cacheResult[T]{
						Err: err.Error(),
						//ErrType: "mongo.ErrNoDocuments",
					}, nil
				}
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
		return result.Data, err
	}

	//If we cached an error, return it
	if result.Err != "" {
		switch result.Err {
		case mongo.ErrNoDocuments.Error():
			return result.Data, mongo.ErrNoDocuments
		}

		return result.Data, errors.New(result.Err)
	}

	//if result.Err != "" {
	//	if knownErr, exists := errorMapping[result.ErrType]; exists {
	//		return result.Data, knownErr
	//	}
	//	return result.Data, errors.New(result.Err)
	//}

	return result.Data, nil
}
