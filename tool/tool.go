package tool

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/afret0/wheel/frame"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
)

func GetEnv() string {
	env := os.Getenv("environment")
	if env == "" {
		env = os.Getenv("ENV")
	}
	return env
}

func IsProEnv() bool {
	return GetEnv() == "pro"
}

func Milliseconds() int64 {
	return time.Now().UnixMilli()
}

func MergeConfig(config1 *viper.Viper, config2 *viper.Viper) *viper.Viper {
	config3 := viper.New()
	for _, key := range config1.AllKeys() {
		config3.Set(key, config1.Get(key))
	}
	for _, key := range config2.AllKeys() {
		config3.Set(key, config2.Get(key))
	}

	return config3
}

func ConStringToInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	return i, err
}

func ConStringToInt64WithoutErr(s string) int64 {
	i, _ := ConStringToInt64(s)
	return i
}

func UUIDWithoutHyphen() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func OpId(ctx context.Context) string {
	opIdValue := ctx.Value("opId")
	opId, ok := opIdValue.(string)
	if !ok {
		return UUIDWithoutHyphen()
	}
	return opId
}

func GrpcCtx(ctx context.Context) context.Context {
	opId := OpId(ctx)

	//md := metadata.Pairs("opid", opId)
	//
	//if md, ok := metadata.FromIncomingContext(ctx); ok {
	//	if val, exists := md["opid"]; exists && len(val) > 0 {
	//		opId = val[0]
	//	} else {
	//		md["opid"] = []string{opId}
	//		ctx = metadata.NewOutgoingContext(ctx, md)
	//	}
	//}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.Pairs("opid", opId)
		ctx = metadata.NewOutgoingContext(ctx, md)
	} else {
		if val, exists := md["opid"]; exists && len(val) > 0 {
			opId = val[0]
		} else {
			md["opid"] = []string{opId}
			//newMd := metadata.Join(md, metadata.Pairs("opid", opId))
			//ctx = metadata.NewOutgoingContext(ctx, newMd)
			ctx = metadata.NewOutgoingContext(ctx, md)
		}
	}

	return ctx
}

func OpIdWithoutDefault(ctx context.Context) string {
	opIdValue := ctx.Value("opId")
	opId, ok := opIdValue.(string)
	if !ok {
		return ""
	}
	return opId
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func NewCtxBK() context.Context {
	return context.WithValue(context.Background(), "opId", strings.ReplaceAll(uuid.New().String(), "-", ""))
}

func RenewCtx(ctx context.Context) context.Context {
	opId := OpId(ctx)
	return context.WithValue(context.Background(), "opId", opId)
}

func ClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	oip := ctx.GetHeader("X-Original-Forwarded-For")

	if oip != "" {
		ip = oip
	}

	return ip
}

func HostId() string {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = UUIDWithoutHyphen()
	}

	return hostname
}

func BoolPtr(b bool) *bool {
	return &b
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func MD5(s string) string {
	if s == "" {
		return ""
	}

	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func UrlContainsCMS(ctx context.Context) bool {
	reqUrl := frame.Request(ctx).RequestURI
	if strings.Contains(reqUrl, "CMS") {
		return true
	}
	return false
}
