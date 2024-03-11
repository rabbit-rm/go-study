package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("rabbit.rm")
		if err != nil {
			fmt.Println("cookie not found")
		}
		c.SetCookie("rabbit.rm", "rabbit.rm99@gmail.com", 60*60*32, "/", "http://192.168.1.233/", false, true)
		fmt.Println(cookie)
	})
	log.Fatal(e.Run(":80"))
}
