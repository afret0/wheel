package tool

import (
	"fmt"
	"reflect"
	"strings"
)

func f(format string, a ...interface{}) string {
	placeholders := make([]string, len(a))
	for i := range placeholders {
		placeholders[i] = "%v"
	}
	format = strings.Replace(format, "{}", "%v", -1)
	for i, v := range a {
		if reflect.TypeOf(v).Kind() == reflect.String {
			placeholders[i] = "%s"
		}
	}
	return fmt.Sprintf(strings.NewReplacer(
		"{", "%",
		"}", "",
	).Replace(format), a...)
}
