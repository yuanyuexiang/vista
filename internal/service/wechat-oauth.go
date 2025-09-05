package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"vista/config"
	"vista/internal/model"
	"vista/internal/repository"
	"vista/pkg/logger"
)

// WechatAuthResult 微信授权结果
type WechatAuthResult struct {
	OpenID       string `json:"openid"`
	UnionID      string `json:"unionid,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// WechatUserInfo 微信用户信息
type WechatUserInfo struct {
	OpenID     string        `json:"openid"`
	UnionID    string        `json:"unionid,omitempty"`
	Nickname   string        `json:"nickname"`
	Sex        int           `json:"sex"`
	Language   string        `json:"language"`
	Province   string        `json:"province"`
	City       string        `json:"city"`
	Country    string        `json:"country"`
	HeadImgURL string        `json:"headimgurl"`
	Privilege  []interface{} `json:"privilege"` // 微信返回的是数组
}

// WechatOAuthService 微信 OAuth 服务
type WechatOAuthService struct {
	userRepo *repository.WechatUserRepository
}

// NewWechatOAuthService 创建微信 OAuth 服务实例
func NewWechatOAuthService(userRepo *repository.WechatUserRepository) *WechatOAuthService {
	return &WechatOAuthService{
		userRepo: userRepo,
	}
}

// ExchangeCodeForToken 用授权码换取 access_token
func (s *WechatOAuthService) ExchangeCodeForToken(code string) (*WechatAuthResult, error) {
	cfg := config.Get()

	// 构建请求 URL
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		cfg.Wechat.AppID, cfg.Wechat.AppSecret, code)

	logger.Infof("Requesting access token from WeChat API - AppID: %s, Code: %s", cfg.Wechat.AppID, code)

	// 发送 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		logger.Errorf("Failed to request access token from WeChat API: %v", err)
		return nil, fmt.Errorf("failed to request access token: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to read response body from WeChat API: %v", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	logger.Infof("WeChat API response for access token: %s", string(body))

	// 解析响应
	var result WechatAuthResult
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Errorf("Failed to parse WeChat API response: %v, Body: %s", err, string(body))
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// 检查是否有错误
	if result.OpenID == "" {
		// 尝试解析错误响应
		var errorResp struct {
			ErrCode int    `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.ErrCode != 0 {
			logger.Errorf("WeChat API returned error: %d - %s", errorResp.ErrCode, errorResp.ErrMsg)
			return nil, fmt.Errorf("wechat api error: %d - %s", errorResp.ErrCode, errorResp.ErrMsg)
		}
		logger.Errorf("Invalid response from WeChat API: %s", string(body))
		return nil, fmt.Errorf("invalid response from wechat api")
	}

	logger.Infof("Successfully got access token - OpenID: %s, ExpiresIn: %d", result.OpenID, result.ExpiresIn)
	return &result, nil
}

// GetUserInfo 获取微信用户信息
func (s *WechatOAuthService) GetUserInfo(accessToken, openID string) (*WechatUserInfo, error) {
	// 构建请求 URL
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openID)

	logger.Infof("Requesting user info from WeChat API - OpenID: %s", openID)

	// 发送 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		logger.Errorf("Failed to request user info from WeChat API: %v", err)
		return nil, fmt.Errorf("failed to request user info: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to read user info response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	logger.Infof("WeChat API user info response: %s", string(body))

	// 解析响应
	var userInfo WechatUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		logger.Errorf("Failed to parse user info response: %v, Body: %s", err, string(body))
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	// 检查是否有错误
	if userInfo.OpenID == "" {
		var errorResp struct {
			ErrCode int    `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.ErrCode != 0 {
			logger.Errorf("WeChat API user info error: %d - %s", errorResp.ErrCode, errorResp.ErrMsg)
			return nil, fmt.Errorf("wechat api error: %d - %s", errorResp.ErrCode, errorResp.ErrMsg)
		}
		logger.Errorf("Invalid user info response from WeChat API: %s", string(body))
		return nil, fmt.Errorf("invalid user info response from wechat api")
	}

	logger.Infof("Successfully got user info - OpenID: %s, Nickname: %s", userInfo.OpenID, userInfo.Nickname)
	return &userInfo, nil
}

// SaveUserAuth 保存用户授权信息到数据库
func (s *WechatOAuthService) SaveUserAuth(authResult *WechatAuthResult, userInfo *WechatUserInfo) error {
	// 计算令牌过期时间
	expiresAt := time.Now().Add(time.Duration(authResult.ExpiresIn) * time.Second)

	// 将 Privilege 数组转换为 JSON 字符串
	privilegeJSON := "[]"
	if len(userInfo.Privilege) > 0 {
		if privilegeBytes, err := json.Marshal(userInfo.Privilege); err == nil {
			privilegeJSON = string(privilegeBytes)
		}
	}

	// 创建数据库用户模型
	user := &model.WechatUser{
		OpenID:       authResult.OpenID,
		UnionID:      authResult.UnionID,
		Nickname:     userInfo.Nickname,
		HeadImgURL:   userInfo.HeadImgURL,
		Sex:          userInfo.Sex,
		Language:     userInfo.Language,
		Province:     userInfo.Province,
		City:         userInfo.City,
		Country:      userInfo.Country,
		Privilege:    privilegeJSON,
		AccessToken:  authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
		ExpiresAt:    expiresAt,
		Scope:        authResult.Scope,
		IsActive:     true,
	}

	// 保存或更新用户信息
	if err := s.userRepo.Save(user); err != nil {
		logger.Errorf("Failed to save user to database: %v", err)
		return fmt.Errorf("failed to save user to database: %v", err)
	}

	logger.Infof("Successfully saved user auth - OpenID: %s, Nickname: %s", userInfo.OpenID, userInfo.Nickname)
	return nil
}

// GetUserByOpenID 根据 openid 获取用户信息
func (s *WechatOAuthService) GetUserByOpenID(openid string) (*model.WechatUser, error) {
	return s.userRepo.GetByOpenID(openid)
}
