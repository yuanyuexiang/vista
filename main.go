package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vista/config"
	"vista/database"
	"vista/internal/router"
	"vista/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize config: %v", err))
	}

	// 初始化日志
	if err := logger.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// 设置Gin模式
	gin.SetMode(config.C.Server.Mode)

	// 初始化数据库
	if err := database.Init(); err != nil {
		logger.Logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// 自动迁移数据库
	// if err := database.AutoMigrate(); err != nil {
	// 	logger.Logger.Fatal("Failed to migrate database", zap.Error(err))
	// }

	// 创建Gin引擎
	r := gin.New()

	// 设置路由
	router.Setup(r)

	// 启动HTTP服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.C.Server.Port),
		Handler:        r,
		ReadTimeout:    config.C.Server.ReadTimeout,
		WriteTimeout:   config.C.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 在goroutine中启动服务器
	go func() {
		logger.Logger.Info("Server starting",
			zap.Int("port", config.C.Server.Port),
			zap.String("mode", config.C.Server.Mode),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Logger.Info("Shutting down server...")

	// 设置5秒的超时时间用于现有连接的完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// 关闭数据库连接
	if err := database.Close(); err != nil {
		logger.Logger.Error("Failed to close database", zap.Error(err))
	}

	logger.Logger.Info("Server exited gracefully")
}
