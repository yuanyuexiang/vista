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
	"vista/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize config: %v", err))
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	// 创建 GraphQL handler
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// 设置路由
	http.Handle("/", playground.Handler("GraphQL playground", "/wechat/query"))
	http.Handle("/wechat/query", srv)

	// 启动HTTP服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.C.Server.Port),
		ReadTimeout:    config.C.Server.ReadTimeout,
		WriteTimeout:   config.C.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 在goroutine中启动服务器
	go func() {
		fmt.Printf("Server starting on port %d\n", config.C.Server.Port)
		fmt.Printf("GraphQL playground: http://localhost:%d/\n", config.C.Server.Port)
		fmt.Printf("GraphQL endpoint: http://localhost:%d/wechat/query\n", config.C.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	// 设置5秒的超时时间用于现有连接的完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	// 关闭数据库连接
	if err := database.Close(); err != nil {
		fmt.Printf("Failed to close database: %v\n", err)
	}

	fmt.Println("Server exited gracefully")
}
