package router

import (
	"github.com/gin-gonic/gin"
	"sample/service"
	"sample/source/middleware"
)

func RegisterRouter(e *gin.Engine) {
	svr := service.GetService()
	router := e.Group("/sample")
	router.Use(middleware.AuthMiddleware())
	router.GET("/sample", svr.Test)
}
