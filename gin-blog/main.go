package main

import (
	"blog/internal/logger"
	"blog/internal/server"
	_ "blog/internal/server/db"

	"go.uber.org/zap"
)

func main() {
	if err := server.Initialized(); err != nil {
		logger.L().Fatal("initializing server failed", zap.Error(err))
	}
}
