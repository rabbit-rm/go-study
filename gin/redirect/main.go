package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.New()
	e.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})
	e.POST("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/foo")
	})
	e.POST("/test2", func(c *gin.Context) {
		// 路由重定向
		c.Request.URL.Path = "/foo"
		c.Request.Method = "POST"
		e.HandleContext(c)
	})
	e.GET("/foo", func(c *gin.Context) {
		c.String(200, "get foo")
	})
	e.POST("/foo", func(c *gin.Context) {
		c.String(200, "post foo")
	})
	log.Fatal(e.Run(":80"))
}
