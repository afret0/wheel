package main

import "github.com/gin-gonic/gin"

func Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	router := gin.Default()
	router.GET("/ping", Ping)

	router.Run(":8080")
}
