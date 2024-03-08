package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Login struct {
	UserName string `json:"username" xml:"username" form:"username" binding:"required"`
	Password string `json:"password" xml:"password" form:"password" binding:"required"`
}

func main() {
	e := gin.New()
	e.Use(gin.Recovery())

	/*
		{
		    "username":"rabbit.rm",
		    "password":"root@123456"
		}
	*/
	e.POST("/bindJson", func(c *gin.Context) {
		var user Login
		if err := c.MustBindWith(&user, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user.UserName != "rabbit.rm" || user.Password != "root@123456" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username|password error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	/*
		<root>
		    <username>rabbit.rm</username>
		    <password>root@123456</password>
		</root>
	*/
	e.POST("/bindXML", func(c *gin.Context) {
		var user Login
		if err := c.MustBindWith(&user, binding.XML); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user.UserName != "rabbit.rm" || user.Password != "root@123456" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username|password error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	e.POST("/bindForm", func(c *gin.Context) {
		var user Login
		if err := c.MustBindWith(&user, binding.Form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user.UserName != "rabbit.rm" || user.Password != "root@123456" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username|password error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	log.Fatal(e.Run(":80"))
}
