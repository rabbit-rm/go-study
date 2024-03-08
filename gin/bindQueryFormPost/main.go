package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name     string `form:"name"`
	Address  string `form:"address"`
	Birthday string `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func main() {
	engine := gin.Default()
	engine.GET("/test", func(c *gin.Context) {
		var p Person
		if err := c.ShouldBind(&p); err == nil {
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
	log.Fatal(engine.Run(":80"))
}
