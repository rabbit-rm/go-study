package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog/internal/config"
	"blog/internal/logger"
	"blog/internal/server/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Initialized() error {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	router.InitRouter(engine)
	s := &http.Server{
		Addr:    config.Host(),
		Handler: engine,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logger.L().Fatal("server error", zap.Error(err))
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	for s := range sigChan {
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			logger.L().Info("service exit signal", zap.Any("signal", s))
		}
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	if err := s.Shutdown(timeout); err != nil {
		logger.L().Fatal("server shutdown error", zap.Error(err))
	}
	logger.L().Info("server exiting")
	return nil
}
