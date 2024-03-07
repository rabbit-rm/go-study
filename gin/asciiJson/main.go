package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.GET("/asciiJson", func(context *gin.Context) {
		context.AsciiJSON(http.StatusOK, gin.H{
			"lang": "Go语言",
			"tag":  "<br>",
		})
	})
	_ = engine.Run(":80")
}
