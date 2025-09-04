package router

import (
	"time"
	"vista/internal/handler"
	"vista/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	// 添加全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.Timeout(30 * time.Second))

	// 初始化处理器
	oauthHandler := handler.NewWechatAuthHandler()

	// API路由组
	api := r.Group("/api/v1")
	{
		// OAuth令牌管理
		oauth := api.Group("/oauth")
		{
			oauth.GET("/authorize", oauthHandler.GetAuthURL) // 获取授权URL
			//oauth.GET("/tokens", oauthHandler.GetTokens)     // 获取令牌列表
		}

	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "vista",
		})
	})

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Relay Server",
			"version": "1.0.0",
			"docs":    "/docs",
		})
	})
}
