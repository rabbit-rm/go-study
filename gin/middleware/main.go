package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 新建一个没有任何中间件的 Engine
	gin.SetMode("release")
	e := gin.New()
	// 定义全局中间件
	// recover 任何 panic，如果产生 panic,会写入500
	e.Use(gin.Recovery())
	// 将日志写入 gin.DefaultWriter，即使你将 GIN_MODE 设置为 release
	e.Use(gin.Logger())
	e.Use(func(c *gin.Context) {
		fmt.Println("Global middleware 1")
		c.Next()
	})
	// 可以为每个路由添加任意中间件
	e.GET("/benchmark", func(c *gin.Context) {
		fmt.Println("benchmark middleware 1")
		c.Next()
	}, func(c *gin.Context) {
		fmt.Println("benchmark middleware 2")
		c.Next()
	}, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "benchmark success",
		})
	})
	group := e.Group("/", func(c *gin.Context) {
		fmt.Println("group(/) Global middleware 1")
		c.Next()
	})
	{
		group.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "login success",
			})
		})
		group.POST("/user", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "add user success",
			})
		})
		// 嵌套路由组
		admin := group.Group("admin", func(c *gin.Context) {
			fmt.Println("group(/admin) Global middleware 1")
		})
		{
			admin.GET("/all", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"status": "all success",
				})
			})
		}
	}

	group2 := e.Group("/", func(c *gin.Context) {
		fmt.Println("group2(/) Global middleware 1")
		c.Next()
	})
	{
		group2.POST("/login2", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "login2 success",
			})
		})
		group2.POST("/user2", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "add user2 success",
			})
		})
		// 嵌套路由组
		admin := group2.Group("admin2", func(c *gin.Context) {
			fmt.Println("group2(/admin2) Global middleware 1")
		})
		{
			admin.GET("/all2", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"status": "all2 success",
				})
			})
		}
	}
	e.Run(":80")
}
