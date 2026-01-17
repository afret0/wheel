package safePointTool

import (
	"math"
	"math/bits"
)

func SafeAddPoints(oldPoints, delta int64) (int64, bool) {
	lo, carry := bits.Add64(uint64(oldPoints), uint64(delta), 0)

	// 如果 carry != 0，说明溢出
	if carry != 0 {
		return 0, false
	}

	// lo 是否能放进 int64？
	if lo > math.MaxInt64 {
		return 0, false
	}

	return int64(lo), true
}

func SafeCalcPoints(giftCount, giftValue int64) (int64, bool) {
	if giftCount < 0 || giftValue < 0 {
		return 0, false
	}

	lo, hi := bits.Mul64(uint64(giftCount), uint64(giftValue))
	if hi != 0 || lo > math.MaxInt64 {
		return 0, false
	}

	return int64(lo), true
}
