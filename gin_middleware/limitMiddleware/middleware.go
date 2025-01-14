package limitMiddleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"reflect"
	"sort"
	"time"
)

type LimitMiddleware struct {
	redis   redis.UniversalClient
	option  *Option
	limiter *redis_rate.Limiter
}

type Option struct {
	Prefix                string
	Rate                  int
	Duration              time.Duration
	NoUidDistinguishable  bool // 是否不区分 uid  默认区分
	HeaderDistinguishable bool // 是否区分 header 默认不区分
}

func PerDuration(rate int, duration time.Duration) redis_rate.Limit {
	return redis_rate.Limit{
		Rate:   rate,
		Burst:  rate,
		Period: duration / time.Duration(rate),
	}
}

func New(rdb redis.UniversalClient, opt *Option) *LimitMiddleware {
	if opt.Prefix == "" {
		panic("prefix is required")
	}

	return &LimitMiddleware{
		redis:   rdb,
		option:  opt,
		limiter: redis_rate.NewLimiter(rdb),
	}
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

func (l *LimitMiddleware) calculateHeaderMD5(ctx *gin.Context) (string, error) {
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

func (l *LimitMiddleware) calculateBodyMD5(ctx *gin.Context) (string, error) {
	// 读取 body
	body, err := ctx.GetRawData()
	if err != nil {
		return "", err
	}

	// 将 body 写回，以便后续中间件和处理函数使用
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 如果 body 为空，返回空字符串
	if len(body) == 0 {
		return "", nil
	}

	// 计算 MD5
	hasher := md5.New()
	hasher.Write(body)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (l *LimitMiddleware) LimitMiddleware(optChain ...*Option) gin.HandlerFunc {
	opt := l.option
	if len(optChain) > 0 {
		opt = optChain[0]
		val := reflect.ValueOf(*opt)
		typeOfConfig := val.Type()

		for i := 0; i < val.NumField(); i++ {
			if isZero(val.Field(i)) {
				reflect.ValueOf(opt).Elem().FieldByName(typeOfConfig.Field(i).Name).Set(reflect.ValueOf(l.option).Elem().Field(i))
			}
		}
	}

	return func(ctx *gin.Context) {
		var err error
		headerS := ""

		if opt.HeaderDistinguishable {
			headerS, err = l.calculateHeaderMD5(ctx)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "header marshal error"})
				ctx.Abort()
				return
			}
		}

		uid := ""
		if !opt.NoUidDistinguishable {
			uid = ctx.Request.Header.Get("_uid")
		}

		//k := fmt.Sprintf("%s:limit:middleware:%s:%s:%s:%s", opt.Prefix, ctx.Request.RequestURI, ctx.Request.Method, uid, headerS)
		//if ctx.Request.Method == http.MethodPost {
		bodyMD5, err := l.calculateBodyMD5(ctx)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "body read error"})
			return
		}

		k := fmt.Sprintf("%s:limit:middleware:%s:%s:%s:%s:%s", opt.Prefix, ctx.Request.RequestURI, ctx.Request.Method, uid, headerS, bodyMD5)
		//}

		allowRet, err := l.limiter.Allow(ctx, k, PerDuration(opt.Rate, opt.Duration))
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": err.Error()})
			ctx.Abort()
			return
		}
		if allowRet.Allowed < 1 {
			ctx.JSON(http.StatusOK, gin.H{"code": -2, "message": "休息 休息 休息一下~"})
			ctx.Abort()
			return
		}
	}
}
