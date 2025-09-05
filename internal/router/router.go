package router

import (
	"vista/internal/handler"
	"vista/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 静态文件服务（用于测试）
	r.StaticFile("/wechat/test", "./test_wechat.html")
	r.StaticFile("/wechat/demo", "./frontend_demo.html") // 前端集成演示页面

	// 微信授权相关路由
	wechat := r.Group("/wechat")
	{
		wechat.GET("/auth", handler.WechatAuth)         // 发起微信授权
		wechat.GET("/callback", handler.WechatCallback) // 微信授权回调
	}

	// 健康检查
	r.GET("/health", handler.HealthCheck)
}
