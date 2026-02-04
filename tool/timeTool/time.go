package timeTool

import (
	"fmt"
	"time"
)

func Week() string {
	year, week := LocalNow().ISOWeek()
	return fmt.Sprintf("%d%02d", year, week)
}

func WeekByTime(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d%02d", year, week)
}

func LastWeek() string {
	// 获取7天前的时间
	lastWeek := LocalNow().AddDate(0, 0, -7)
	year, week := lastWeek.ISOWeek()
	return fmt.Sprintf("%d%02d", year, week)
}

func WeekDay() time.Weekday {
	return LocalNow().Weekday()
}

func Year() string {
	return fmt.Sprintf("%d", LocalNow().Year())
	//return fmt.Sprintf("%d", LocalNow().Year())
}

// LastMonth 返回上个月的年月格式（例如：202312）
func LastMonth() string {
	now := LocalNow()

	// 获取当前年月
	currentYear := now.Year()
	currentMonth := now.Month()

	// 计算上个月的年月
	var year int
	var month time.Month

	if currentMonth == time.January {
		year = currentYear - 1
		month = time.December
	} else {
		year = currentYear
		month = currentMonth - 1
	}

	return fmt.Sprintf("%d%02d", year, month)
}

func Month() string {
	return fmt.Sprintf("%s%02d", Year(), LocalNow().Month())
}

func Day() string {
	return fmt.Sprintf("%s%02d", Month(), LocalNow().Day())
}

func Hour() string {
	return fmt.Sprintf("%s%02d", Day(), LocalNow().Hour())
}

func Minute() string {
	return fmt.Sprintf("%s%02d", Hour(), LocalNow().Minute())
}

func Second() string {
	return fmt.Sprintf("%s%02d", Minute(), LocalNow().Second())
}

func MidnightTody() time.Time {
	return MidnightToday()
}

func MidnightToday() time.Time {
	//location, err := time.LoadLocation("Asia/Shanghai")
	//if err != nil {
	//	panic(err)
	//}
	//
	//return time.Date(LocalNow().Year(), LocalNow().Month(), LocalNow().Day(), 0, 0, 0, 0, location)

	//location := time.FixedZone("UTC+8", 8*60*60) // 东八区，偏移量为8小时(8*60*60秒)

	now := LocalNow()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, TZBeijing())
}

// Deprecated: use TZBeijing instead
func Location() *time.Location {
	//return Location()
	location := time.FixedZone("UTC+8", 8*60*60)
	return location
}

func TZBeijing() *time.Location {
	location := time.FixedZone("UTC+8", 8*60*60)
	return location
}

func LocalNow() time.Time {
	return Now()
}

func Now() time.Time {
	return time.Now().In(TZBeijing())
}

func ParseMillisecond(ts int64) time.Time {
	t := time.Unix(0, ts*int64(time.Millisecond))
	return t.In(TZBeijing())
}

func ParseSecond(ts int64) time.Time {
	t := time.Unix(ts, 0)
	return t.In(TZBeijing())
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func ParseTimeStr(timeStr string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, TZBeijing())
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
