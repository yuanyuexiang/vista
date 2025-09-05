package handler

import (
	"fmt"
	"net/http"
	"time"
	"vista/internal/database"
	"vista/internal/repository"
	"vista/internal/service"
	"vista/pkg/response"

	"github.com/gin-gonic/gin"
)

// WechatAuthRequest 前端发送的授权请求
type WechatAuthRequest struct {
	Code  string `json:"code" binding:"required"` // 微信授权码
	State string `json:"state,omitempty"`         // 可选的state参数
}

// WechatAuthByCode 通过授权码获取用户信息
func WechatAuthByCode(c *gin.Context) {
	var req WechatAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	fmt.Printf("Received auth request - Code: %s, State: %s\n", req.Code, req.State)

	// 创建服务实例
	userRepo := repository.NewWechatUserRepository(database.GetDB())
	wechatService := service.NewWechatOAuthService(userRepo)

	fmt.Println("Starting WeChat OAuth process...")

	// 使用 code 换取 access_token 和用户信息
	authResult, err := wechatService.ExchangeCodeForToken(req.Code)
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

	// 构建返回的用户信息
	userInfoData := gin.H{
		"openid":     authResult.OpenID,
		"nickname":   userInfo.Nickname,
		"headimgurl": userInfo.HeadImgURL,
		"sex":        userInfo.Sex,
		"language":   userInfo.Language,
		"country":    userInfo.Country,
		"province":   userInfo.Province,
		"city":       userInfo.City,
		"privilege":  userInfo.Privilege,
		"login_time": time.Now().Unix(),
	}

	// 返回JSON响应
	response.Success(c, gin.H{
		"message":   "微信登录成功",
		"user_info": userInfoData,
	})
}

// GetUserInfo 获取用户信息，包含授权状态检查
func GetUserInfo(c *gin.Context) {
	// 从URL参数获取openid
	openid := c.Param("openid")
	if openid == "" {
		// 也可以从query参数获取
		openid = c.Query("openid")
	}

	if openid == "" {
		response.BadRequest(c, "openid is required")
		return
	}

	// 从数据库获取用户信息
	repo := repository.NewWechatUserRepository(database.GetDB())
	user, err := repo.GetByOpenID(openid)
	if err != nil {
		// 用户不存在，需要重新授权
		response.Success(c, gin.H{
			"exists":    false,
			"need_auth": true,
			"message":   "用户不存在，需要重新授权",
		})
		return
	}

	// 检查授权是否过期（24小时）
	now := time.Now()
	expireTime := user.UpdatedAt.Add(24 * time.Hour)
	needReauth := now.After(expireTime)

	// 返回用户信息和授权状态
	response.Success(c, gin.H{
		"exists":    true,
		"need_auth": needReauth,
		"user_info": gin.H{
			"openid":     user.OpenID,
			"nickname":   user.Nickname,
			"sex":        user.Sex,
			"language":   user.Language,
			"city":       user.City,
			"province":   user.Province,
			"country":    user.Country,
			"headimgurl": user.HeadImgURL,
			"unionid":    user.UnionID,
			"privilege":  user.Privilege,
			"is_active":  user.IsActive,
			"login_time": user.UpdatedAt.Unix(),
			"expires_at": expireTime.Unix(),
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

	response.Success(c, gin.H{
		"status":   "ok",
		"service":  "vista-wechat-auth",
		"database": "connected",
	})
}
