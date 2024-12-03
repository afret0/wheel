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

}
