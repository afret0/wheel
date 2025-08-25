package tool

import "testing"

func Test_RandInt64InRange(t *testing.T) {
	m := map[int64]int64{}
	min := int64(1)
	max := int64(10)
	for i := 0; i < 1000000; i++ {
		m[RandInt64InRange(min, max)] += 1
	}
	t.Logf("m: %+v", m)
}

func Test_FormatToWan(t *testing.T) {
	l := []int64{1, 2, 300000, -1, -100, -2 - 23234234, -234234}

	for _, num := range l {
		t.Logf("num: %d, to wan: %s", num, FormatToWan(num))
	}
}
