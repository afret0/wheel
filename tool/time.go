package tool

import (
	"fmt"
	"time"
)

func Week() string {
	year, week := time.Now().ISOWeek()
	return fmt.Sprintf("%d%d", year, week)
}

func WeekDay() time.Weekday {
	return time.Now().Weekday()
}

func Year() string {
	return fmt.Sprintf("%d", time.Now().Year())
}

func Month() string {
	return fmt.Sprintf("%s%d", Year(), time.Now().Month())
}

func Day() string {
	return fmt.Sprintf("%s%d", Month(), time.Now().Day())
}

func Hour() string {
	return fmt.Sprintf("%s%d", Day(), time.Now().Hour())
}

func Minute() string {
	return fmt.Sprintf("%s%d", Hour(), time.Now().Minute())
}

func Second() string {
	return fmt.Sprintf("%s%d", Minute(), time.Now().Second())
}

func MidnightTody() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
}
