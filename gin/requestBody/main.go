package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.New()
	e.POST("/bindBody", bindBodyTo)
	e.POST("/bindBodyWith", bindBodyWithTo)
	e.Run(":80")
}
