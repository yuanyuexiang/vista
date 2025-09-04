package handler

import (
	"time"
	"vista/internal/dto"
	"vista/internal/service"
	"vista/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type WechatAuthHandler struct {
	wechatService *service.WechatAuthService
	validator     *validator.Validate
}

func NewWechatAuthHandler() *WechatAuthHandler {
	return &WechatAuthHandler{
		wechatService: service.NewWechatAuthService(),
		validator:     validator.New(),
	}
}

// 微信公众号授权回调处理
func (h *WechatAuthHandler) AuthCallback(c *gin.Context) {
	var req dto.WechatAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithDetails(c, "invalid request body", err.Error())
		return
	}
	if err := h.validator.Struct(&req); err != nil {
		response.BadRequestWithDetails(c, "validation failed", err.Error())
		return
	}
	// 这里应调用微信API获取 access_token、openid 等信息
	// 示例：
	wechatResp := &dto.WechatAuthResponse{
		OpenID:      "mock_openid",
		AccessToken: "mock_access_token",
		ExpiresIn:   7200,
		Scope:       "snsapi_userinfo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// 保存授权信息
	if err := h.wechatService.SaveWechatAuth(wechatResp); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "微信授权成功", wechatResp)
}
func (h *WechatAuthHandler) GetAuthURL(c *gin.Context) {
	// 这里应生成微信授权URL
	authURL := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect"
	response.Success(c, gin.H{"url": authURL})
}
