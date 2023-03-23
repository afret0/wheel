package main

import (
	"fmt"
	"strings"
)

func a() {
	fmt.Println(1)
}

func b() {
	defer a()
	fmt.Println(2)
}

func c() {

	fmt.Println(3)
}

func main() {
	m := make(map[string]string, 0)
	m["1"] = "a"
	m["2"] = "b"
	m["3"] = "c"
	l := make([]string, 0)
	for k, v := range m {
		s := fmt.Sprintf("%s=%s", k, v)
		fmt.Println(s)
		l = append(l, s)
	}
	s := strings.Join(l,"&")
	fmt.Println(l,s)
}
