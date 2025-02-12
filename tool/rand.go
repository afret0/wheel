package tool

import (
	"math/rand"
	"time"
)

// RandInt64InRange 生成 [min, max) 的 int64 随机数
func RandInt64InRange(min, max int64) int64 {
	if min >= max {
		panic("max must be greater than min")
	}

	rand.NewSource(time.Now().UnixNano())
	return min + rand.Int63n(max-min)
}
