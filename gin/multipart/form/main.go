package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	e.POST("/user", func(c *gin.Context) {
		// 从表单中获取数据
		msg := c.PostForm("message")
		// 从表单中获取数据，如果没有就返回默认值
		nick := c.DefaultPostForm("nick", "anonymous")
		c.JSON(200, gin.H{
			"status":  "success",
			"message": msg,
			"nick":    nick,
		})
	})
	e.Run(":80")
}
