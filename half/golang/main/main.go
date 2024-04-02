package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func g1() {
	fmt.Printf("g1 start\n")
	time.Sleep(10 * time.Second)
	fmt.Printf("g1 end\n")
}

func g2() {
	fmt.Printf("g2 start\n")
	go g1()
	time.Sleep(5 * time.Second)
	fmt.Printf("g2 end\n")
}

func main() {
	g2()
	time.Sleep(3*time.Second)
	// router := gin.Default()
	// router.GET("/ping", Ping)

	// router.Run(":8080")
}
