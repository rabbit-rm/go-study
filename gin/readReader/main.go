package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.New()
	// 使用中间件
	e.Use(gin.Recovery())
	e.GET("download", func(c *gin.Context) {
		resp, err := http.Get("http://192.168.80.51:10002/demo/d/public/wp.svg")
		if err != nil || resp.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		body := resp.Body
		length := resp.ContentLength
		contentType := resp.Header.Get("Content-Type")
		c.DataFromReader(200, length, contentType, body, map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		})
	})
	e.Run(":80")

}
