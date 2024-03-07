package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	_ = engine.Run("192.168.1.233:80")
}
