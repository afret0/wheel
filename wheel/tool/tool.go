package tool

import "fmt"

func SprintfStruct(s interface{}) string {
	return fmt.Sprintf("%+v", s)
}
