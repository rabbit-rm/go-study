package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	root = filepath.Join(root, "staticFile", "resources")
	e := gin.Default()
	e.Static("/assets", filepath.Join(root, "assets"))
	e.StaticFS("/ico", http.Dir(filepath.Join(root, "ico")))
	e.StaticFile("/hello.html", filepath.Join(root, "hello.html"))
	log.Fatal(e.Run(":80"))
}
