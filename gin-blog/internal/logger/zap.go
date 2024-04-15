package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLogger *zap.Logger

func L() *zap.Logger {
	return zapLogger
}

func init() {
	var err error
	cfg := zap.NewProductionConfig()
	// time encoder
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05")
	zapLogger, err = cfg.Build(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic(err)
	}
}
