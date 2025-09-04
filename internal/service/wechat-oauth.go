package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"vista/config"
	"vista/database"
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
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid,omitempty"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgURL string `json:"headimgurl"`
}

// WechatOAuthService 微信 OAuth 服务
type WechatOAuthService struct{}

// ExchangeCodeForToken 用授权码换取 access_token
func (s *WechatOAuthService) ExchangeCodeForToken(code string) (*WechatAuthResult, error) {
	// 构建请求 URL
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		config.C.Wechat.AppID, config.C.Wechat.AppSecret, code)

	// 发送 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request access token: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 解析响应
	var result WechatAuthResult
	if err := json.Unmarshal(body, &result); err != nil {
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
			return nil, fmt.Errorf("wechat api error: %d - %s", errorResp.ErrCode, errorResp.ErrMsg)
		}
		return nil, fmt.Errorf("invalid response from wechat api")
	}

	return &result, nil
}

// GetUserInfo 获取微信用户信息
func (s *WechatOAuthService) GetUserInfo(accessToken, openID string) (*WechatUserInfo, error) {
	// 构建请求 URL
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openID)

	// 发送 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request user info: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 解析响应
	var userInfo WechatUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	// 检查是否有错误
	if userInfo.OpenID == "" {
		var errorResp struct {
			ErrCode int    `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.ErrCode != 0 {
			return nil, fmt.Errorf("wechat api error: %d - %s", errorResp.ErrCode, errorResp.ErrMsg)
		}
		return nil, fmt.Errorf("invalid user info response from wechat api")
	}

	return &userInfo, nil
}

// SaveUserAuth 保存用户授权信息到数据库
func (s *WechatOAuthService) SaveUserAuth(authResult *WechatAuthResult, userInfo *WechatUserInfo) error {
	// 计算令牌过期时间
	expiresAt := time.Now().Add(time.Duration(authResult.ExpiresIn) * time.Second)

	// 创建数据库用户模型
	user := &database.WechatUser{
		OpenID:       authResult.OpenID,
		UnionID:      authResult.UnionID,
		Nickname:     userInfo.Nickname,
		Avatar:       userInfo.HeadImgURL,
		Sex:          userInfo.Sex,
		Province:     userInfo.Province,
		City:         userInfo.City,
		Country:      userInfo.Country,
		AccessToken:  authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
		ExpiresAt:    expiresAt,
		Scope:        authResult.Scope,
		IsActive:     true,
	}

	// 保存或更新用户信息
	if err := database.SaveOrUpdateWechatUser(user); err != nil {
		return fmt.Errorf("failed to save user to database: %v", err)
	}

	fmt.Printf("成功保存用户授权信息 - OpenID: %s, Nickname: %s\n", userInfo.OpenID, userInfo.Nickname)
	return nil
}

// GetUserByOpenID 根据 openid 获取用户信息
func (s *WechatOAuthService) GetUserByOpenID(openid string) (*database.WechatUser, error) {
	return database.GetWechatUserByOpenID(openid)
}
