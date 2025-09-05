package model

import "time"

// WechatUser 微信用户模型
type WechatUser struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	OpenID     string `json:"openid" gorm:"column:openid;type:varchar(64);unique;not null;comment:微信用户唯一标识"`
	UnionID    string `json:"unionid" gorm:"column:unionid;type:varchar(64);comment:微信开放平台统一ID"`
	Nickname   string `json:"nickname" gorm:"column:nickname;type:varchar(128);comment:用户昵称"`
	Sex        int    `json:"sex" gorm:"column:sex;type:tinyint;default:0;comment:性别 0未知 1男 2女"`
	Language   string `json:"language" gorm:"column:language;type:varchar(16);comment:语言"`
	Country    string `json:"country" gorm:"column:country;type:varchar(32);comment:国家"`
	Province   string `json:"province" gorm:"column:province;type:varchar(32);comment:省份"`
	City       string `json:"city" gorm:"column:city;type:varchar(32);comment:城市"`
	HeadImgURL string `json:"headimgurl" gorm:"column:headimgurl;type:varchar(512);comment:头像URL"`
	Privilege  string `json:"privilege" gorm:"column:privilege;type:json;comment:用户特权信息(JSON数组)"`

	// 授权相关字段
	AccessToken  string    `json:"access_token" gorm:"column:access_token;type:varchar(512);comment:访问令牌"`
	RefreshToken string    `json:"refresh_token" gorm:"column:refresh_token;type:varchar(512);comment:刷新令牌"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"column:expires_at;comment:令牌过期时间"`
	Scope        string    `json:"scope" gorm:"column:scope;type:varchar(64);comment:授权作用域"`
	IsActive     bool      `json:"is_active" gorm:"column:is_active;default:true;comment:用户是否活跃"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (WechatUser) TableName() string {
	return "wechat_users"
}
