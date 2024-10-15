package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/afret0/wheel/lock"
	"github.com/afret0/wheel/log"
)

type RepositoryCache struct {
	*Repository
	cache redis.UniversalClient
	opt   *RepositoryCacheOption
	lock  *lock.Locker
	debug bool
}

var ErrCacheMiss = errors.New("cache miss")

func (r *RepositoryCache) GenCacheK(filter interface{}) (string, error) {
	fs, err := json.Marshal(filter)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:dbcache:%s:%s", r.opt.Prefix, r.Repository.collection.Name(), fs), nil
}

type RepositoryCacheOption struct {
	Prefix string
	TTL    int
	Debug  bool
}

type RCO = RepositoryCacheOption

func GetRepositoryCache(repo *Repository, cache redis.UniversalClient, opt *RepositoryCacheOption) *RepositoryCache {

	if opt.Prefix == "" {
		panic("prefix is empty")
	}

	if opt.TTL == 0 {
		opt.TTL = 60 * 60
	}

	r := new(RepositoryCache)
	r.Repository = repo
	r.cache = cache
	r.opt = opt
	r.lock = lock.GetLocker(cache)
	r.debug = opt.Debug

	return r
}

//type RepositoryCacheFindOneOptions struct {
//	TTL int
//}

func (r *RepositoryCache) getFromCache(ctx context.Context, entity interface{}, key string) error {
	//lg := log.CtxLogger(ctx)

	ret, _ := r.cache.Get(ctx, key).Result()
	if ret == "" {
		return ErrCacheMiss
	}

	if ret == mongo.ErrNoDocuments.Error() {
		return mongo.ErrNoDocuments
	}

	err := json.Unmarshal([]byte(ret), entity)
	if err != nil {
		//lg.Errorf("json unmarshal failed, err: %v", err)
		return err
	}

	if r.debug {
		log.GetLogger().Infof("cache hit, key: %s", key)
	}
	return nil
}

func (r *RepositoryCache) FindOne(ctx context.Context, entity interface{}, filter interface{}, opts ...*RepositoryCacheOption) error {
	//lg := log.CtxLogger(ctx).WithField("filter", filter)
	opt := r.opt
	if len(opts) > 1 {
		opt = opts[0]
	}

	key, err := r.GenCacheK(filter)
	if err != nil {
		//lg.Errorf("gen cache key failed, err: %v", err)
		return err
	}

	err = r.getFromCache(ctx, entity, key)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return mongo.ErrNoDocuments
	}

	if err == nil {
		return nil
	}

	//obtain, err := r.lock.ObtainWaitRetry(ctx, fmt.Sprintf("%s:obtain:%s", r.opt.Prefix, key), 1, 3, 1)
	//if err != nil {
	//	lg.Errorf("obtain failed, err: %v", err)
	//	return err
	//}
	//defer func() {
	//	_ = obtain.Release(ctx)
	//}()
	//
	// err = r.getFromCache(ctx, entity, key)
	// if errors.Is(err, mongo.ErrNoDocuments) {
	// 	return mongo.ErrNoDocuments
	// }
	// if err == nil {
	// 	return nil
	// }

	err = r.Repository.FindOne(ctx, entity, filter)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		//lg.Errorf("find one failed, err: %v", err)
		return err
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		_ = r.cache.Set(ctx, key, err.Error(), time.Duration(opt.TTL)*time.Second).Err()
		//if err != nil {
		//lg.Errorf("cache set failed, err: %v", err)
		//}
		return mongo.ErrNoDocuments
	}

	bs, err := json.Marshal(entity)
	if err != nil {
		//lg.Errorf("json marshal failed, err: %v", err)
		return err
	}

	err = r.cache.Set(ctx, key, bs, time.Duration(opt.TTL)*time.Second).Err()
	if err != nil {
		//lg.Errorf("cache set failed, err: %v", err)
	}

	if r.debug {
		log.GetLogger().Infof("cache miss, key: %s", key)
	}

	return nil
}

func (r *RepositoryCache) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	//lg := log.CtxLogger(ctx).WithField("filter", filter)

	key, err := r.GenCacheK(filter)
	if err != nil {
		//lg.Errorf("gen cache key failed, err: %v", err)
		return nil, err
	}
	defer func() {
		r.DelCache(ctx, key)
	}()

	result, err := r.Repository.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		//lg.Errorf("update one failed, err: %v", err)
		return nil, err
	}

	return result, nil
}

func (r *RepositoryCache) Find(ctx context.Context, entityList interface{}, filter interface{}, opts ...*options.FindOptions) error {
	return fmt.Errorf("多个条件无法刷新缓存, 需要单独处理")
}

func (r *RepositoryCache) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, fmt.Errorf("多个条件无法刷新缓存, 需要单独处理")
}

func (r *RepositoryCache) FindOneAndUpdate(ctx context.Context, entity interface{}, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	//lg := log.CtxLogger(ctx).WithField("filter", filter)

	key, err := r.GenCacheK(filter)
	if err != nil {
		//lg.Errorf("gen cache key failed, err: %v", err)
		return err
	}

	defer func() {
		r.DelCache(ctx, key)
	}()

	one := r.collection.FindOneAndUpdate(ctx, filter, update, opts...)
	err = one.Decode(entity)
	if err != nil {
		//lg.Errorf("find one and update failed, err: %v", err)
		return err
	}

	return nil
}

func (r *RepositoryCache) InsertOne(ctx context.Context, doc interface{}, filter interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	key, err := r.GenCacheK(filter)
	if err != nil {
		//lg.Errorf("gen cache key failed, err: %v", err)
		return nil, err
	}

	defer func() {
		r.DelCache(ctx, key)
	}()

	return r.Repository.InsertOne(ctx, doc, opts...)
}

func (r *RepositoryCache) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return nil, fmt.Errorf("无法刷新缓存, 需要单独处理")
}

func (r *RepositoryCache) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	//lg := log.CtxLogger(ctx).WithField("filter", filter)
	key, err := r.GenCacheK(filter)
	if err != nil {
		//lg.Errorf("gen cache key failed, err: %v", err)
		return nil, err
	}

	defer func() {
		r.DelCache(ctx, key)
	}()

	result, err := r.collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RepositoryCache) DelCache(ctx context.Context, key string) {
	r.cache.Del(ctx, key)
	go func() {
		time.Sleep(1 * time.Second)
		r.cache.Del(context.Background(), key)
	}()
}

func (r *RepositoryCache) DelCacheByFilter(ctx context.Context, filter interface{}) {
	key, err := r.GenCacheK(filter)
	if err != nil {
		return
	}
	r.DelCache(ctx, key)
}
