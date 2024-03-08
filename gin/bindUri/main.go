package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name string `uri:"name" form:"name" binding:"required"`
	ID   string `uri:"id" form:"id" binding:"required,uuid"`
}

func main() {
	e := gin.Default()
	e.GET("/:id/:name", func(c *gin.Context) {
		var p Person
		if err := c.BindUri(&p); err == nil {
			c.JSON(200, gin.H{
				"status": "success",
				"person": p,
			})
		} else {
			c.JSON(400, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
		}
	})
	log.Fatal(e.Run(":80"))
}
