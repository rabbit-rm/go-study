package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func mapParam() {
	e := gin.New()
	e.POST("/params", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")
		c.JSON(200, gin.H{
			"ids":   ids,
			"names": names,
		})
	})
	log.Fatal(e.Run(":80"))
}
