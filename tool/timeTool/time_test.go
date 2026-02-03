package timeTool

import (
	"fmt"
	"testing"
)

func Test_time(t *testing.T) {
	fmt.Printf("LastWeek: %s\n", LastWeek())
	fmt.Printf("Week: %s\n", Week())
	fmt.Printf("WeekDay: %s\n", WeekDay())
	fmt.Printf("Year: %s\n", Year())
	fmt.Printf("LastMonth: %s\n", LastMonth())
	fmt.Printf("Month: %s\n", Month())
	fmt.Printf("Day: %s\n", Day())
	fmt.Printf("Hour: %s\n", Hour())
	fmt.Printf("Minute: %s\n", Minute())
	fmt.Printf("Second: %s\n", Second())

	fmt.Printf("midnight: %s\n", MidnightTody())
}

func Test_ParseMillisecond(t *testing.T) {

	tm := ParseMillisecond(1770103802548)

	t.Logf("ParseMillisecond: %s", FormatTime(tm))
}

func Test_ParseTimeStr(t *testing.T) {
	tmStr := "2026-02-03 15:30:02"

	tm, err := ParseTimeStr(tmStr)
	if err != nil {
		t.Errorf("ParseTimeStr error: %v", err)
		return
	}

	t.Logf("ParseTimeStr: %s", tm)
}
