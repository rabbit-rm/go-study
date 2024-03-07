package main

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	e.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.Abort()
		}
		headers := form.File["file"]
		var headerNames []string
		for _, header := range headers {
			headerNames = append(headerNames, header.Filename)
			_ = c.SaveUploadedFile(header, filepath.Join("W:\\GoProject\\private\\gin\\uploadFile", header.Filename))
		}
		c.JSON(200, gin.H{
			"result":    "success",
			"fileNames": headerNames,
		})
	})
	e.Run(":80")
}

func singleFile(e *gin.Engine) gin.IRoutes {
	return e.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.Abort()
		}
		_ = c.SaveUploadedFile(file, filepath.Join("W:\\GoProject\\private\\gin\\uploadFile", file.Filename))
		c.JSON(200, gin.H{
			"result":   "success",
			"fileName": file.Filename,
		})
	})
}
