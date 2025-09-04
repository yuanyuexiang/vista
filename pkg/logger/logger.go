package logger

import (
	"os"
	"vista/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init() error {
	var cfg zap.Config

	switch config.C.Log.Format {
	case "console":
		cfg = zap.NewDevelopmentConfig()
	default:
		cfg = zap.NewProductionConfig()
	}

	// 设置日志级别
	level := zapcore.InfoLevel
	switch config.C.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(level)

	// 设置输出
	switch config.C.Log.Output {
	case "stderr":
		cfg.OutputPaths = []string{"stderr"}
	case "file":
		cfg.OutputPaths = []string{config.C.Log.FilePath}
		// 确保日志目录存在
		if err := os.MkdirAll("logs", 0755); err != nil {
			return err
		}
	default:
		cfg.OutputPaths = []string{"stdout"}
	}

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		return err
	}

	return nil
}

func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}
