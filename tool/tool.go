package tool

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func GetEnv() string {
	env := os.Getenv("environment")
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

func OpIdWithoutDefault(ctx context.Context) string {
	opIdValue := ctx.Value("opId")
	opId, ok := opIdValue.(string)
	if !ok {
		return ""
	}
	return opId
}

func MergeByJson(from interface{}, to interface{}) {
	fromJson, _ := json.Marshal(from)
	_ = json.Unmarshal(fromJson, to)
}

func SecondsUntilMidnight() int64 {
	now := time.Now()
	// 获取明天的日期
	tomorrow := now.AddDate(0, 0, 1)
	// 创建明天零点的时间
	midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, now.Location())
	// 计算从现在到明天零点的时间差
	duration := midnight.Sub(now)
	// 返回时间差的秒数
	return int64(duration.Seconds())
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
