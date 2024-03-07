package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	e.SecureJsonPrefix("Rabbit.RM")
	e.GET("/sjson", func(c *gin.Context) {
		c.SecureJSON(200, gin.H{
			"result": []string{"user", "foo", "bar"},
		})
	})
	e.Run(":80")
}
