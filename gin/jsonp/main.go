package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	e.GET("/jsonp", func(c *gin.Context) {
		c.JSONP(200, map[string]string{
			"foo": "bar",
		})
	})
	e.Run(":80")
}
