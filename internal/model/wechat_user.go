package model

import "time"

// WechatUser 微信用户模型
type WechatUser struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	OpenID       string    `json:"openid" gorm:"type:varchar(128);uniqueIndex;not null;comment:微信openid"`
	UnionID      string    `json:"unionid" gorm:"type:varchar(128);index;comment:微信unionid"`
	Nickname     string    `json:"nickname" gorm:"type:varchar(100);comment:用户昵称"`
	Avatar       string    `json:"avatar" gorm:"type:text;comment:头像URL"`
	Sex          int       `json:"sex" gorm:"type:tinyint;comment:性别 0-未知 1-男 2-女"`
	Province     string    `json:"province" gorm:"type:varchar(50);comment:省份"`
	City         string    `json:"city" gorm:"type:varchar(50);comment:城市"`
	Country      string    `json:"country" gorm:"type:varchar(50);comment:国家"`
	AccessToken  string    `json:"access_token" gorm:"type:text;comment:访问令牌"`
	RefreshToken string    `json:"refresh_token" gorm:"type:text;comment:刷新令牌"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"comment:令牌过期时间"`
	Scope        string    `json:"scope" gorm:"type:varchar(100);comment:授权范围"`
	IsActive     bool      `json:"is_active" gorm:"default:true;comment:是否有效"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (WechatUser) TableName() string {
	return "wechat_users"
}
