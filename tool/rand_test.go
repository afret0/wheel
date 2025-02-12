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
