package tool

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/status"

	"github.com/afret0/wheel/frame"
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
