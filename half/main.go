package main

import "fmt"

func main() {
	s1 := "asdf"
	s2 := "qwer"

	switch 1 {
	case len(s1):
		fmt.Println(s1)
	case len(s2):
		fmt.Println(s2)
	}
}
