package cacheInterceptor

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"reflect"
	"time"
)

type Option struct {
	TTL                     int64
	MetadataDistinguishable bool
}

type CacheInterceptor struct {
	prefix        string
	redis         redis.UniversalClient
	opt           *Option
	MethodToCache map[string]*Option
}

func New(redisClient redis.UniversalClient, prefix string, methodToCache map[string]*Option, opt *Option) *CacheInterceptor {
	if prefix == "" {
		panic("prefix is required")
	}

	return &CacheInterceptor{
		prefix:        prefix,
		redis:         redisClient,
		opt:           opt,
		MethodToCache: methodToCache,
	}
}

func (c *CacheInterceptor) AddMethodToCache(method string, opt *Option) {
	if opt == nil {
		opt = c.opt
	}

	c.MethodToCache[method] = opt
}

type cacheStruc struct {
	res any
	err error
}

func (c *CacheInterceptor) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if opt, ok := c.MethodToCache[info.FullMethod]; ok {
		lg := CtxLogger(ctx).WithField("method", info.FullMethod)

		val := reflect.ValueOf(*opt)
		typeOfOpt := val.Type()

		for i := 0; i < val.NumField(); i++ {
			if isZero(val.Field(i)) {
				reflect.ValueOf(typeOfOpt).Elem().FieldByName(typeOfOpt.Field(i).Name).Set(reflect.ValueOf(c.opt).Elem().Field(i))
			}
		}

		reqB, err := json.Marshal(req)
		if err != nil {
			lg.Errorf("序列化请求失败, err: %+v", err)
			return nil, err
		}

		var mdB []byte
		if opt.MetadataDistinguishable {
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				mdB, err = json.Marshal(md)
				if err != nil {
					lg.Errorf("序列化元数据失败, err: %+v", err)
					return nil, err
				}
			}
		}
		hasher := md5.New()
		hasher.Write(mdB)
		mdMD5 := hex.EncodeToString(hasher.Sum(nil))

		k := fmt.Sprintf("%s:%s:%s:%s", c.prefix, info.FullMethod, mdMD5, string(reqB))

		cacheResS, err := c.redis.Get(ctx, k).Result()
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				lg.Errorf("读取缓存失败, err: %+v", err)
			} else {
				lg.Info("缓存未命中")
			}
		}

		if cacheResS != "" {
			var cache cacheStruc
			err = json.Unmarshal([]byte(cacheResS), &cache)
			if err != nil {
				lg.Errorf("反序列化缓存失败, err: %+v", err)
				return nil, err
			}

			return cache.res, cache.err
		}

		res, err := handler(ctx, req)

		cache := cacheStruc{
			res: res,
			err: err,
		}

		cacheResB, err := json.Marshal(cache)
		if err != nil {
			lg.Errorf("序列化缓存失败, err: %+v", err)
			return nil, err
		}

		err = c.redis.Set(ctx, k, cacheResB, time.Duration(opt.TTL)*time.Second).Err()
		if err != nil {
			lg.Errorf("写入缓存失败, err: %+v", err)
		}

		return res, err
	}

	return handler(ctx, req)
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
