package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		_, _ = fmt.Fprintf(gin.DefaultWriter, "endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	e := gin.New()
	e.Use(gin.Logger())
	out, err := os.Create("W:\\GoProject\\private\\study\\gin\\logger\\gin.log")
	if err != nil {
		fmt.Printf("error:%+v\n", err)
		os.Exit(1)
	}
	gin.DefaultWriter = io.MultiWriter(out, os.Stdout)
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})
	e.Run(":80")
}
