package main

import (
	"log"
	"vista/config"
	"vista/database"
	"vista/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatal("Failed to initialize config:", err)
	}

	// 初始化日志
	if err := logger.Init(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		logger.Logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 自动迁移
	// if err := database.AutoMigrate(); err != nil {
	// 	logger.Logger.Fatal("Failed to migrate database", zap.Error(err))
	// }

	logger.Logger.Info("Database seeded successfully!")
}
