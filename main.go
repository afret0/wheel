package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func GetServiceName() string {
	fmt.Println("please input service name...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	sentence := scanner.Text()
	fmt.Printf("new service: %s", sentence)
	return sentence
}

func RegisterRouter(service string) {
	data := []byte("hello")
	file, err := os.Create("./router/test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func NewService(service string) {

}

func main() {
	//serviceName := GetServiceName()
	//if serviceName == "" {
	//	fmt.Println("The service name you entered is empty. Exit")
	//	log.Fatal()
	//}
	RegisterRouter("123")
}
