package service

import (
	"fmt"
	"time"
	"vista/graph/model"
)

type WechatAuthService struct{}

func NewWechatAuthService() *WechatAuthService {
	return &WechatAuthService{}
}

// SaveWechatAuth 保存微信授权信息
func (s *WechatAuthService) SaveWechatAuth(resp *model.WechatAuth) error {
	//db := database.GetDB()
	// 这里只是示例，实际应有微信用户表
	// err := db.Create(resp).Error
	// return err
	fmt.Println("保存微信授权信息：", resp)
	return nil
}

// GetWechatAuthByOpenID 根据 openid 获取微信授权信息
func (s *WechatAuthService) GetWechatAuthByOpenID(openid string) (*model.WechatAuth, error) {
	//db := database.GetDB()
	// 示例：实际应有微信用户表
	// var user WechatUser
	// err := db.Where("openid = ?", openid).First(&user).Error
	// if err != nil {
	//     return nil, err
	// }
	// return &model.WechatAuth{...}, nil
	return &model.WechatAuth{
		Openid:      openid,
		AccessToken: "mock_access_token_" + openid,
		ExpiresIn:   7200,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}, nil
}
