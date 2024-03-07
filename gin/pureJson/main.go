package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	e.GET("/asciiJson", func(c *gin.Context) {
		c.AsciiJSON(http.StatusOK, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	e.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	e.GET("/pureJson", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	e.Run(":80")
}
