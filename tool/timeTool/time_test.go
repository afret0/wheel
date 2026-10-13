package timeTool

import (
	"fmt"
	"testing"
	"time"
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

func Test_ByTime(t *testing.T) {
	bj := TZBeijing()

	tests := []struct {
		name          string
		at            time.Time
		wantYear      string
		wantLastMonth string
		wantMonth     string
		wantDay       string
		wantHour      string
		wantMinute    string
		wantSecond    string
	}{
		{
			name:          "普通时间",
			at:            time.Date(2026, time.March, 5, 8, 9, 7, 0, bj),
			wantYear:      "2026",
			wantLastMonth: "202602",
			wantMonth:     "202603",
			wantDay:       "20260305",
			wantHour:      "2026030508",
			wantMinute:    "202603050809",
			wantSecond:    "20260305080907",
		},
		{
			name:          "年初月初零点",
			at:            time.Date(2026, time.January, 1, 0, 0, 0, 0, bj),
			wantYear:      "2026",
			wantLastMonth: "202512",
			wantMonth:     "202601",
			wantDay:       "20260101",
			wantHour:      "2026010100",
			wantMinute:    "202601010000",
			wantSecond:    "20260101000000",
		},
		{
			name:          "月末最后一纳秒",
			at:            time.Date(2026, time.March, 31, 23, 59, 59, int(time.Second)-1, bj),
			wantYear:      "2026",
			wantLastMonth: "202602",
			wantMonth:     "202603",
			wantDay:       "20260331",
			wantHour:      "2026033123",
			wantMinute:    "202603312359",
			wantSecond:    "20260331235959",
		},
		{
			name:          "31 号减一个月不溢出到下月",
			at:            time.Date(2026, time.May, 31, 12, 0, 0, 0, bj),
			wantYear:      "2026",
			wantLastMonth: "202604",
			wantMonth:     "202605",
			wantDay:       "20260531",
			wantHour:      "2026053112",
			wantMinute:    "202605311200",
			wantSecond:    "20260531120000",
		},
		{
			name:          "按入参自带时区格式化, 不做隐式转换",
			at:            time.Date(2026, time.February, 28, 20, 30, 15, 0, time.UTC),
			wantYear:      "2026",
			wantLastMonth: "202601",
			wantMonth:     "202602",
			wantDay:       "20260228",
			wantHour:      "2026022820",
			wantMinute:    "202602282030",
			wantSecond:    "20260228203015",
		},
		{
			name:          "毫秒时间戳反解后截断到秒",
			at:            ParseMillisecond(time.Date(2026, time.October, 13, 17, 31, 0, int(time.Millisecond)*265, bj).UnixMilli()),
			wantYear:      "2026",
			wantLastMonth: "202609",
			wantMonth:     "202610",
			wantDay:       "20261013",
			wantHour:      "2026101317",
			wantMinute:    "202610131731",
			wantSecond:    "20261013173100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := []struct {
				name string
				got  string
				want string
			}{
				{"YearByTime", YearByTime(tt.at), tt.wantYear},
				{"LastMonthByTime", LastMonthByTime(tt.at), tt.wantLastMonth},
				{"MonthByTime", MonthByTime(tt.at), tt.wantMonth},
				{"DayByTime", DayByTime(tt.at), tt.wantDay},
				{"HourByTime", HourByTime(tt.at), tt.wantHour},
				{"MinuteByTime", MinuteByTime(tt.at), tt.wantMinute},
				{"SecondByTime", SecondByTime(tt.at), tt.wantSecond},
			}

			for _, g := range got {
				if g.got != g.want {
					t.Errorf("%s(%s) = %s, want %s", g.name, FormatTime(tt.at), g.got, g.want)
				}
			}
		})
	}
}

func Test_noArgFuncEqualsByTime(t *testing.T) {
	before := LocalNow()

	levels := []struct {
		name string
		got  string
		by   func(time.Time) string
	}{
		{name: "Year", got: Year(), by: YearByTime},
		{name: "LastMonth", got: LastMonth(), by: LastMonthByTime},
		{name: "Month", got: Month(), by: MonthByTime},
		{name: "Day", got: Day(), by: DayByTime},
		{name: "Hour", got: Hour(), by: HourByTime},
		{name: "Minute", got: Minute(), by: MinuteByTime},
		{name: "Second", got: Second(), by: SecondByTime},
	}

	after := LocalNow()

	for _, lv := range levels {
		lo, hi := lv.by(before), lv.by(after)

		if lv.got < lo || lv.got > hi {
			t.Errorf("%s() = %s, want within [%s, %s]", lv.name, lv.got, lo, hi)
		}
	}
}
