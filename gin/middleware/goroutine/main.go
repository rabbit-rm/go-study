package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/long_async", func(c *gin.Context) {
		cc := c.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			fmt.Println("Done! in path " + cc.Request.URL.Path)
		}()
	})
	e.GET("/long_sync", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		fmt.Println("Done! in path " + c.Request.URL.Path)
	})
	e.Run(":80")
}
