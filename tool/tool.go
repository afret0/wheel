package tool

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/status"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: false, TimestampFormat: "2006-01-02 15:04:05"})
}

// GetEnv deprecated, use Env instead
func GetEnv() string {
	return Env()
}

func Env() string {
	env := os.Getenv("environment")
	if env == "" {
		env = os.Getenv("ENV")
	}
	return env
}

func IsProEnv() bool {
	return GetEnv() == "pro"
}

func IsDevEnv() bool {
	return GetEnv() == "dev"
}

func IsTestEnv() bool {
	return GetEnv() == "test"
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

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// 网关(如 ingress)会把最外层客户端 IP 放进 X-Original-Forwarded-For,
// 其余代理链依次追加到 X-Forwarded-For, 因此前者优先级最高。
const (
	headerOriginalForwardedFor = "X-Original-Forwarded-For"
	headerForwardedFor         = "X-Forwarded-For"
	headerRealIP               = "X-Real-Ip"
)

// ClientIP 返回请求方的真实 IP, 行为与旧版 ClientIP(*gin.Context) 一致。
// ctx 直接传 *gin.Context, 或以它为父级派生出来的 context 都可以。
//
// 能取到 *gin.Context 时复用 gin 的 ClientIP(), 保留 engine 的 TrustedPlatform /
// TrustedProxies / RemoteIPHeaders 配置语义; 只能取到 *http.Request 时按
// X-Original-Forwarded-For > X-Forwarded-For > X-Real-Ip > RemoteAddr 的顺序解析。
// 两条路径都优先返回网关注入的 X-Original-Forwarded-For。
//
// 注意: c.Request.Context() 不携带 gin 写入的 key, 这种情况下取不到 IP(返回空串),
// 需要 client IP 就直接把 c 往下传。
func ClientIP(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if gc, ok := ctx.Value(gin.ContextKey).(*gin.Context); ok && gc != nil && gc.Request != nil {
		if oip := gc.Request.Header.Get(headerOriginalForwardedFor); oip != "" {
			return oip
		}

		return gc.ClientIP()
	}

	return clientIPFromRequest(requestFromCtx(ctx))
}

// clientIPFromRequest 在没有 gin.Context(engine 配置)可用时, 从 *http.Request 上尽力解析客户端 IP。
func clientIPFromRequest(req *http.Request) string {
	if req == nil {
		return ""
	}

	if oip := req.Header.Get(headerOriginalForwardedFor); oip != "" {
		return oip
	}

	if xff := req.Header.Get(headerForwardedFor); xff != "" {
		ip, _, _ := strings.Cut(xff, ",")
		if ip = strings.TrimSpace(ip); ip != "" {
			return ip
		}
	}

	if rip := req.Header.Get(headerRealIP); rip != "" {
		return rip
	}

	remoteAddr := strings.TrimSpace(req.RemoteAddr)
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}

	return ip
}

// requestFromCtx 取出 ctx 上的 *http.Request, 与 frame.Request 等价。
// tool 不能直接调 frame.Request: frame -> database -> tool 已经存在依赖,
// 再加 tool -> frame 就成了 import cycle, 因此这里就地实现同样的查找。
func requestFromCtx(ctx context.Context) *http.Request {
	req, ok := ctx.Value(gin.ContextRequestKey).(*http.Request)
	if !ok {
		return nil
	}

	return req
}

func HostId() string {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = UUIDWithoutHyphen()
	}

	return hostname
}

var (
	appNameFallbackOnce sync.Once
	osHostnameOnce      = sync.OnceValue(readOsHostname)
)

func AppName() string {
	hostId := HostId()
	if !isGeneratedUuid(hostId) {
		return appNameFromHostId(hostId)
	}

	hostname := osHostnameOnce()
	if hostname == "" {
		return ""
	}

	appNameFallbackOnce.Do(func() {
		logrus.Infof("appName fallback, reason: HOSTNAME 为空使 HostId 退化成随机 uuid, 改用 os.Hostname() 取服务名, hostname: %s", hostname)
	})

	return appNameFromHostId(hostname)
}

func readOsHostname() string {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		logrus.Warnf("osHostname fail, reason: os.Hostname() 取不到主机名, 服务名留空不写 caller, err: %v", err)

		return ""
	}

	return hostname
}

func appNameFromHostId(hostId string) string {
	parts := strings.Split(hostId, "-")
	if len(parts) < 3 {
		return hostId
	}

	return strings.Join(parts[:len(parts)-2], "-")
}

func isGeneratedUuid(s string) bool {
	if len(s) != 32 {
		return false
	}

	for _, c := range s {
		if !(c >= '0' && c <= '9') && !(c >= 'a' && c <= 'f') {
			return false
		}
	}

	return true
}

// Deprecated
func BoolPtr(b bool) *bool {
	return &b
}

// Deprecated
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

func FormatToWan(num int64) string {
	fuhao := ""
	if num < 0 {
		fuhao = "-"
	}
	num = int64(math.Abs(float64(num)))

	if num < 10000 {
		return fmt.Sprintf("%s%d", fuhao, num)
	}

	wan := float64(num) / 10000.0
	return fmt.Sprintf("%s%.2fw", fuhao, wan)
}

func IsInterfaceNil(v interface{}) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

func ErrCode(err error) int {

	st, ok := status.FromError(err)
	if !ok {
		return -1
	}

	code := st.Code()
	return int(code)
}

func Debug(keyChains ...string) bool {
	key := "DEBUG"
	if len(keyChains) > 0 {
		key = keyChains[0]
	}

	return EnvEnabled(key)
}

func EnvEnabled(key string) bool {
	tk := os.Getenv(key)

	switch tk {
	case "true", "TRUE", "1", "yes", "YES":
		return true
	}

	return false
}

const Boy = 1
const Girl = 2

func GenderFromID(id string) (int, error) {
	if len(id) == 18 {
		sexCode := id[16] - '0'
		if sexCode%2 == 0 {
			return Girl, nil
		}
		return Boy, nil
	}

	if len(id) == 15 {
		sexCode := id[14] - '0'
		if sexCode%2 == 0 {
			return Girl, nil
		}
		return Boy, nil
	}

	return 0, fmt.Errorf("id err")
}

func FindFieldByJSONTag(v any, tag string) (fieldName string, ok bool) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("json") == tag {
			return f.Name, true
		}
	}
	return "", false
}

func FindFieldByBSONTag(v any, tag string) (fieldName string, ok bool) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("bson") == tag {
			return f.Name, true
		}
	}
	return "", false
}

func ExtractFieldValueByBSONTag(v any, tag string) (any, bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, false
	}

	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.Tag.Get("bson") == tag {
			fv := rv.Field(i)
			if fv.IsValid() {
				return fv.Interface(), true
			}
		}
	}
	return nil, false
}

func IsK8sHayo() bool {
	return os.Getenv("K8S-APP") == "hayo"
}
