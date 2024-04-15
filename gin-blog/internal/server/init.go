package server

import (
	"blog/internal/config"
	"blog/internal/server/router"

	"github.com/gin-gonic/gin"
)

func Initialized() error {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	router.InitRouter(engine)
	return engine.Run(config.Host())
}
