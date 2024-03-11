package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.New()
	e.GET("/user/:name", func(c *gin.Context) {
		user := c.Param("name")
		c.JSON(200, gin.H{
			"user": user,
		})
	})
	e.GET("/user/:name/*action", func(c *gin.Context) {
		user := c.Param("name")
		action := c.Param("action")
		c.JSON(200, gin.H{
			"user":   user,
			"action": action,
		})
	})
	log.Fatal(e.Run(":80"))
}
