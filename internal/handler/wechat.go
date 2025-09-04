package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"vista/config"
	"vista/database"
	"vista/internal/service"

	"github.com/gin-gonic/gin"
)

// WechatAuth 发起微信授权，重定向到微信授权页面
func WechatAuth(c *gin.Context) {
	// 生成 state 参数防止 CSRF 攻击
	state := generateState()

	// 构建微信授权 URL
	authURL := buildWechatAuthURL(state)

	// 重定向到微信授权页面
	c.Redirect(http.StatusFound, authURL)
}

// WechatCallback 处理微信授权回调
func WechatCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "authorization code is required",
		})
		return
	}

	// 验证 state 参数（这里简化了，实际项目中应该存储在 session 或 cache 中）
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid state parameter",
		})
		return
	}

	// 使用 code 换取 access_token 和用户信息
	wechatService := &service.WechatOAuthService{}
	authResult, err := wechatService.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to exchange code: %v", err),
		})
		return
	}

	// 获取用户信息
	userInfo, err := wechatService.GetUserInfo(authResult.AccessToken, authResult.OpenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get user info: %v", err),
		})
		return
	}

	// 保存用户授权信息
	if err := wechatService.SaveUserAuth(authResult, userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to save user auth: %v", err),
		})
		return
	}

	// 生成应用自己的 JWT token（这里简化为直接返回微信信息）
	// 实际项目中应该：
	// 1. 根据 openid 查询或创建用户
	// 2. 生成 JWT token
	// 3. 重定向到前端页面，带上 token

	// 这里简化为重定向到前端，带上用户信息
	frontendURL := fmt.Sprintf("http://localhost:3000/auth/success?openid=%s&nickname=%s",
		authResult.OpenID, url.QueryEscape(userInfo.Nickname))

	c.Redirect(http.StatusFound, frontendURL)
}

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	// 检查数据库连接
	if err := database.HealthCheck(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   "error",
			"service":  "vista-wechat-auth",
			"database": "disconnected",
			"error":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"service":  "vista-wechat-auth",
		"database": "connected",
	})
}

// buildWechatAuthURL 构建微信授权 URL
func buildWechatAuthURL(state string) string {
	baseURL := "https://open.weixin.qq.com/connect/oauth2/authorize"
	params := url.Values{}
	params.Add("appid", config.C.Wechat.AppID)
	params.Add("redirect_uri", config.C.Wechat.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "snsapi_userinfo") // 或 snsapi_base
	params.Add("state", state)

	return fmt.Sprintf("%s?%s#wechat_redirect", baseURL, params.Encode())
}

// generateState 生成随机 state 参数
func generateState() string {
	// 这里简化了，实际项目中应该生成更安全的随机字符串
	return "random_state_123"
}
