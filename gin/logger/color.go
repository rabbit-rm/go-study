package main

import (
	"github.com/gin-gonic/gin"
)

func loggerColor() {
	// 禁用日志颜色化
	gin.DisableConsoleColor()
	// 强制日志颜色化
	gin.ForceConsoleColor()
	e := gin.New()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.GET("/data", func(c *gin.Context) {
		c.JSON(200, "success")
	})
	e.Run(":80")
}
