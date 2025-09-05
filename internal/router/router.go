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

	// 演示页面
	r.StaticFile("/", "./wechat_frontend_driven.html")                       // 主页：前端驱动模式
	r.StaticFile("/wechat/frontend-driven", "./wechat_frontend_driven.html") // 前端驱动模式演示页面

	// 核心API接口
	api := r.Group("/api")
	{
		// 通过code获取用户信息并保存
		api.POST("/wechat/auth", handler.WechatAuthByCode)

		// 查询用户信息，判断是否需要重新授权
		api.GET("/user/:openid", handler.GetUserInfo)
	}

	// 健康检查
	r.GET("/health", handler.HealthCheck)
}
