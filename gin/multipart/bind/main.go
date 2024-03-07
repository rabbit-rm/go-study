package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func main() {
	e := gin.Default()
	e.POST("/login", func(c *gin.Context) {
		var from = LoginForm{}
		// 表单绑定结构体
		if err := c.ShouldBind(&from); err == nil {
			if from.User == "admin" && from.Password == "123" {
				c.JSON(http.StatusOK, gin.H{
					"result": "login success",
				})
			} else {
				c.JSON(401, gin.H{
					"result": "login failed",
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprintf("error:%+v", err),
			})
		}
	})
	_ = e.Run(":80")
}
