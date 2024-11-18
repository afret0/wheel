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
