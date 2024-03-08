package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.Use(customLogger())
	engine.Use(gin.Recovery())
	engine.Use(customExample)
	engine.GET("/ping", func(c *gin.Context) {
		time.Sleep(200 * time.Millisecond)
		c.JSON(200, "pong")
	})
	log.Fatal(engine.Run(":80"))
}

// 自定义中间件
func customExample(c *gin.Context) {
	t := time.Now()
	c.Set("example", "rabbit.rm")
	c.Next()
	latency := time.Since(t)
	log.Printf("[%d]%s", c.Writer.Status(), latency)
}

func customLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		// 你的自定义格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC1123),
			params.Method,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)
	})
}
