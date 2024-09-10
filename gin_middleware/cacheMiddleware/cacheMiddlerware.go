package cacheMiddleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"reflect"
	"sort"
	"time"
)

type CacheMiddleware struct {
	redis  redis.UniversalClient
	config *Config
}

type Config struct {
	Prefix                string
	TTL                   int64
	HeaderDistinguishable bool
	NoUidDistinguishable  bool
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

func New(redisClient redis.UniversalClient, config *Config) *CacheMiddleware {
	if config.Prefix == "" {
		panic("prefix is required")
	}

	return &CacheMiddleware{
		redis:  redisClient,
		config: config,
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (cm *CacheMiddleware) calculateHeaderMD5(ctx *gin.Context) (string, error) {
	keys := make([]string, 0, len(ctx.Request.Header))
	for k := range ctx.Request.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedHeader := make(map[string][]string)
	for _, k := range keys {
		sortedHeader[k] = ctx.Request.Header[k]
	}

	headerBytes, err := json.Marshal(sortedHeader)
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	hasher.Write(headerBytes)
	md5 := hex.EncodeToString(hasher.Sum(nil))

	return md5, nil
}

func (cm *CacheMiddleware) CacheMiddleware(configChain ...*Config) gin.HandlerFunc {
	config := cm.config
	if len(configChain) > 0 {
		config = configChain[0]
		val := reflect.ValueOf(*config)
		typeOfConfig := val.Type()

		for i := 0; i < val.NumField(); i++ {
			if isZero(val.Field(i)) {
				reflect.ValueOf(config).Elem().FieldByName(typeOfConfig.Field(i).Name).Set(reflect.ValueOf(cm.config).Elem().Field(i))
			}
		}
	}

	return func(ctx *gin.Context) {
		var err error
		headerS := ""

		if config.HeaderDistinguishable {
			headerS, err = cm.calculateHeaderMD5(ctx)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "header marshal error"})
				return
			}
		}

		uid := ""
		if !config.NoUidDistinguishable {
			uid = ctx.Request.Header.Get("_uid")
		}

		k := fmt.Sprintf("%s:cache:middleware:%s:%s:%s", config.Prefix, ctx.Request.RequestURI, uid, headerS)

		val, err := cm.redis.Get(ctx, k).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			ctx.Next()
			return
		}

		if errors.Is(err, redis.Nil) {
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
			ctx.Writer = blw

			ctx.Next()

			bs := blw.body.String()

			_ = cm.redis.Set(ctx, k, bs, time.Duration(config.TTL)*time.Second)

			return
		}

		n := new(interface{})

		err = json.Unmarshal([]byte(val), n)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": err.Error()})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, n)
		ctx.Abort()
		return
	}
}
