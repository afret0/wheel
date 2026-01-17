package safePointTool

import (
	"math"
)

func SafeAddPoints(oldPoints, delta int64) (int64, bool) {
	if delta > 0 && oldPoints > math.MaxInt64-delta {
		return 0, false
	}
	if delta < 0 && oldPoints < math.MinInt64-delta {
		return 0, false
	}
	return oldPoints + delta, true
}

func SafeCalcPoints(giftCount, giftValue int64) (int64, bool) {
	if giftCount < 0 || giftValue < 0 {
		return 0, false
	}

	// giftCount * giftValue 是否会溢出？
	if giftCount != 0 && giftValue > math.MaxInt64/giftCount {
		return 0, false
	}

	return giftCount * giftValue, true
}
