package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"vista/config"
	"vista/internal/database"
	"vista/internal/repository"
	"vista/internal/service"
	"vista/pkg/response"
	"vista/pkg/utils"

	"github.com/gin-gonic/gin"
)

// WechatAuth 发起微信授权，重定向到微信授权页面
func WechatAuth(c *gin.Context) {
	// 生成 state 参数防止 CSRF 攻击
	state := utils.GenerateState()

	// 构建微信授权 URL
	authURL := buildWechatAuthURL(state)

	fmt.Printf("Generated auth URL: %s\n", authURL)

	// 重定向到微信授权页面
	c.Redirect(http.StatusFound, authURL)
}

// WechatCallback 处理微信授权回调
func WechatCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	fmt.Printf("WeChat callback received - Code: %s, State: %s\n", code, state)

	if code == "" {
		fmt.Println("Error: Missing authorization code")
		response.BadRequest(c, "authorization code is required")
		return
	}

	// 验证 state 参数（这里简化了，实际项目中应该存储在 session 或 cache 中）
	if state == "" {
		fmt.Println("Error: Missing state parameter")
		response.BadRequest(c, "invalid state parameter")
		return
	}

	// 创建服务实例
	userRepo := repository.NewWechatUserRepository(database.GetDB())
	wechatService := service.NewWechatOAuthService(userRepo)

	fmt.Println("Starting WeChat OAuth process...")

	// 使用 code 换取 access_token 和用户信息
	authResult, err := wechatService.ExchangeCodeForToken(code)
	if err != nil {
		fmt.Printf("Error exchanging code for token: %v\n", err)
		response.InternalServerError(c, fmt.Sprintf("微信登录失败: %v", err))
		return
	}

	fmt.Printf("Successfully got access token for OpenID: %s\n", authResult.OpenID)

	// 测试号支持获取完整用户信息
	userInfo, err := wechatService.GetUserInfo(authResult.AccessToken, authResult.OpenID)
	if err != nil {
		fmt.Printf("Error getting user info: %v\n", err)
		response.InternalServerError(c, fmt.Sprintf("获取用户信息失败: %v", err))
		return
	}

	fmt.Printf("Successfully got user info - OpenID: %s, Nickname: %s\n", userInfo.OpenID, userInfo.Nickname)

	// 保存用户授权信息
	if err := wechatService.SaveUserAuth(authResult, userInfo); err != nil {
		fmt.Printf("Error saving user auth: %v\n", err)
		response.InternalServerError(c, fmt.Sprintf("保存用户信息失败: %v", err))
		return
	}

	fmt.Println("Successfully saved user auth to database")

	// 返回成功响应
	response.Success(c, gin.H{
		"message": "微信登录成功",
		"user_info": gin.H{
			"openid":     authResult.OpenID,
			"nickname":   userInfo.Nickname,
			"headimgurl": userInfo.HeadImgURL,
			"sex":        userInfo.Sex,
			"language":   userInfo.Language,
			"country":    userInfo.Country,
			"province":   userInfo.Province,
			"city":       userInfo.City,
			"privilege":  userInfo.Privilege,
			"note":       "使用微信测试号获取完整用户信息",
		},
	})
}

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	// 检查数据库连接
	if err := database.HealthCheck(); err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("database health check failed: %v", err))
		return
	}

	response.SuccessWithMessage(c, "service is healthy", gin.H{
		"status":   "ok",
		"service":  "vista-wechat-auth",
		"database": "connected",
	})
}

// buildWechatAuthURL 构建微信授权 URL
func buildWechatAuthURL(state string) string {
	cfg := config.Get()

	// 使用URL编码的回调地址
	redirectURI := url.QueryEscape(cfg.Wechat.RedirectURI)

	// 构建授权URL - 测试号支持完整功能，使用 snsapi_userinfo
	authURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect",
		cfg.Wechat.AppID,
		redirectURI,
		state,
	)

	fmt.Printf("Building WeChat Test Account auth URL:\n")
	fmt.Printf("  AppID: %s\n", cfg.Wechat.AppID)
	fmt.Printf("  RedirectURI: %s\n", cfg.Wechat.RedirectURI)
	fmt.Printf("  EncodedURI: %s\n", redirectURI)
	fmt.Printf("  State: %s\n", state)
	fmt.Printf("  Final URL: %s\n", authURL)

	return authURL
}
