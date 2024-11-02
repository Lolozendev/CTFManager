package internal

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.SugaredLogger
	once   sync.Once
)

func initLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger = zap.Must(config.Build()).Sugar()
}

func GetLogger() *zap.SugaredLogger {
	once.Do(initLogger)
	return logger
}
