package main

import (
	"github.com/gin-gonic/gin"
)

type person struct {
	Name    string `form:"name"`
	Age     uint8  `form:"age"`
	Country string `form:"country"`
}

func main() {
	mapParam()
}

func bindQuery() {
	e := gin.Default()
	e.GET("/queryBind", func(c *gin.Context) {
		var p person
		err := c.ShouldBindQuery(&p)
		if err != nil {
			c.Abort()
		}
		c.JSON(200, p)
	})
	e.Run(":80")
}

func query() {
	e := gin.Default()
	e.POST("/query", func(c *gin.Context) {
		id := c.DefaultQuery("id", "1")
		name := c.Query("name")
		postId := c.DefaultPostForm("id", "1")
		age := c.DefaultPostForm("age", "23")
		c.JSON(200, gin.H{
			"id":     id,
			"name":   name,
			"age":    age,
			"postId": postId,
		})
	})
	e.Run(":80")
}
