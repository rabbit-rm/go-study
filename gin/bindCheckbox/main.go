package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.New()
	e.LoadHTMLGlob("W:\\GoProject\\private\\study\\gin\\bindCheckbox\\views\\*")
	e.GET("/index", indexHandler)
	e.POST("/index", IndexFormHandler)
	log.Fatal(e.Run(":80"))
}

func indexHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

type colorForm struct {
	Colors []string `form:"colors"`
}

func IndexFormHandler(c *gin.Context) {
	var colors colorForm
	if err := c.ShouldBind(&colors); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "colors": colors})
}
