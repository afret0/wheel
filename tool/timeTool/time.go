package timeTool

import (
	"fmt"
	"strconv"
	"time"
)

func Week() string {
	year, week := time.Now().ISOWeek()
	return fmt.Sprintf("%d%02d", year, week)
}

func LastWeek() string {
	// 获取7天前的时间
	lastWeek := time.Now().AddDate(0, 0, -7)
	year, week := lastWeek.ISOWeek()
	return fmt.Sprintf("%d%02d", year, week)
}

func WeekDay() time.Weekday {
	return time.Now().Weekday()
}

func Year() string {
	return fmt.Sprintf("%d", time.Now().Year())
}

// LastMonth 返回上个月的年月格式（例如：202312）
func LastMonth() string {
	now := time.Now()

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
	return fmt.Sprintf("%s%02d", Year(), time.Now().Month())
}

func Day() string {
	return fmt.Sprintf("%s%02d", Month(), time.Now().Day())
}

func Hour() string {
	return fmt.Sprintf("%s%02d", Day(), time.Now().Hour())
}

func Minute() string {
	return fmt.Sprintf("%s%02d", Hour(), time.Now().Minute())
}

func Second() string {
	return fmt.Sprintf("%s%02d", Minute(), time.Now().Second())
}

func MidnightTody() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
}

// ParseWeek 解析周格式 (例如: 202401)
func ParseWeek(weekStr string) (time.Time, error) {
	if len(weekStr) < 6 {
		return time.Time{}, fmt.Errorf("invalid week format: %s", weekStr)
	}

	year, err := strconv.Atoi(weekStr[:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year in week: %s", weekStr)
	}

	week, err := strconv.Atoi(weekStr[4:])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid week number: %s", weekStr)
	}

	// 获取该年第一天
	jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	// 找到第一周的星期一
	firstMonday := jan1
	if jan1.Weekday() != time.Monday {
		daysUntilMonday := int(time.Monday - jan1.Weekday())
		if daysUntilMonday <= 0 {
			daysUntilMonday += 7
		}
		firstMonday = jan1.AddDate(0, 0, daysUntilMonday)
	}

	// 计算目标周的时间
	target := firstMonday.AddDate(0, 0, (week-1)*7)
	return target, nil
}

// ParseMonth 解析月份格式 (例如: 202401)
func ParseMonth(monthStr string) (time.Time, error) {
	if len(monthStr) < 6 {
		return time.Time{}, fmt.Errorf("invalid month format: %s", monthStr)
	}

	year, err := strconv.Atoi(monthStr[:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year in month: %s", monthStr)
	}

	month, err := strconv.Atoi(monthStr[4:])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month number: %s", monthStr)
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local), nil
}

// ParseDay 解析日期格式 (例如: 20240101)
func ParseDay(dayStr string) (time.Time, error) {
	if len(dayStr) < 8 {
		return time.Time{}, fmt.Errorf("invalid day format: %s", dayStr)
	}

	year, err := strconv.Atoi(dayStr[:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year in day: %s", dayStr)
	}

	month, err := strconv.Atoi(dayStr[4:6])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month in day: %s", dayStr)
	}

	day, err := strconv.Atoi(dayStr[6:])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day number: %s", dayStr)
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

// ParseHour 解析小时格式 (例如: 2024010112)
func ParseHour(hourStr string) (time.Time, error) {
	if len(hourStr) < 10 {
		return time.Time{}, fmt.Errorf("invalid hour format: %s", hourStr)
	}

	t, err := ParseDay(hourStr[:8])
	if err != nil {
		return time.Time{}, err
	}

	hour, err := strconv.Atoi(hourStr[8:])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour number: %s", hourStr)
	}

	return time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, time.Local), nil
}

// ParseMinute 解析分钟格式 (例如: 202401011230)
func ParseMinute(minuteStr string) (time.Time, error) {
	if len(minuteStr) < 12 {
		return time.Time{}, fmt.Errorf("invalid minute format: %s", minuteStr)
	}

	t, err := ParseHour(minuteStr[:10])
	if err != nil {
		return time.Time{}, err
	}

	minute, err := strconv.Atoi(minuteStr[10:])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute number: %s", minuteStr)
	}

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, 0, 0, time.Local), nil
}

// ParseSecond 解析秒格式 (例如: 20240101123045)
func ParseSecond(secondStr string) (time.Time, error) {
	if len(secondStr) < 14 {
		return time.Time{}, fmt.Errorf("invalid second format: %s", secondStr)
	}

	t, err := ParseMinute(secondStr[:12])
	if err != nil {
		return time.Time{}, err
	}

	second, err := strconv.Atoi(secondStr[12:])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid second number: %s", secondStr)
	}

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), second, 0, time.Local), nil
}
