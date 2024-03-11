package gin

import (
	"log"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestExecute(t *testing.T) {
	e := gin.New()
	e.Use(gin.Recovery())
	group := e.Group("/", func(c *gin.Context) {
		c.Set("group(/) middleware 1", "rabbit.rm")
	})
	{
		group.GET("/test", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		group.POST("/post", func(c *gin.Context) {
			c.JSON(200, "post success")
		})
	}
	log.Fatal(e.Run())
}
