package main

import (
	"fmt"
)

func a(){
	fmt.Println(1)
}

func b(){
	defer a()
	fmt.Println(2)
}

func c()  {

	fmt.Println(3)
}

func main() {
	b()
	c()
}
