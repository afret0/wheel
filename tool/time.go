package tool

import (
	"github.com/afret0/wheel/tool/timeTool"
	"time"
)

// Deprecated: 已废弃
func Week() string {
	//year, week := time.Now().ISOWeek()
	//return fmt.Sprintf("%d%d", year, week)
	return timeTool.Week()
}

// Deprecated: 已废弃，请使用 timeTool.LastWeek()
func LastWeek() string {
	// 获取7天前的时间
	return timeTool.LastWeek()
}

// Deprecated: 已废弃
func WeekDay() time.Weekday {
	//return time.Now().Weekday()
	return timeTool.WeekDay()
}

// Deprecated: 已废弃
func Year() string {
	//return fmt.Sprintf("%d", time.Now().Year())
	return timeTool.Year()
}

// Deprecated: 已废弃
func LastMonth() string {
	return timeTool.LastMonth()
}

// Deprecated: 已废弃
func Month() string {
	//return fmt.Sprintf("%s%d", Year(), time.Now().Month())
	return timeTool.Month()
}

// Deprecated: 已废弃
func Day() string {
	//return fmt.Sprintf("%s%d", Month(), time.Now().Day())
	return timeTool.Day()
}

// Deprecated: 已废弃
func Hour() string {
	//return fmt.Sprintf("%s%d", Day(), time.Now().Hour())
	return timeTool.Hour()
}

// Deprecated: 已废弃
func Minute() string {
	//return fmt.Sprintf("%s%d", Hour(), time.Now().Minute())
	return timeTool.Minute()
}

// Deprecated: 已废弃
func Second() string {
	//return fmt.Sprintf("%s%d", Minute(), time.Now().Second())
	return timeTool.Second()
}

// Deprecated: 已废弃
func MidnightTody() time.Time {
	//return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	return timeTool.MidnightTody()
}
