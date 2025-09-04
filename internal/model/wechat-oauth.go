package model

import (
	"time"
)

// OAuthToken OAuth令牌模型
type OAuthToken struct {
	BaseModel
	AccessToken  string    `json:"access_token" gorm:"type:text;not null;comment:访问令牌"`
	TokenType    string    `json:"token_type" gorm:"size:20;default:'Bearer';comment:令牌类型"`
	RefreshToken string    `json:"refresh_token" gorm:"type:text;comment:刷新令牌"`
	Expiry       time.Time `json:"expiry" gorm:"comment:过期时间"`
	ExpiresIn    int64     `json:"expires_in" gorm:"comment:过期秒数"`
	IsActive     bool      `json:"is_active" gorm:"default:true;comment:是否有效"`
}

// OAuthTokenStatus OAuth令牌状态
type OAuthTokenStatus int

const (
	OAuthTokenStatusInactive OAuthTokenStatus = 0 // 无效
	OAuthTokenStatusActive   OAuthTokenStatus = 1 // 有效
	OAuthTokenStatusExpired  OAuthTokenStatus = 2 // 已过期
	OAuthTokenStatusRevoked  OAuthTokenStatus = 3 // 已撤销
)

// String 实现 Stringer 接口
func (s OAuthTokenStatus) String() string {
	switch s {
	case OAuthTokenStatusInactive:
		return "inactive"
	case OAuthTokenStatusActive:
		return "active"
	case OAuthTokenStatusExpired:
		return "expired"
	case OAuthTokenStatusRevoked:
		return "revoked"
	default:
		return "unknown"
	}
}

// TableName 指定表名
func (OAuthToken) TableName() string {
	return "oauth_tokens"
}

// IsExpired 检查令牌是否过期
func (t *OAuthToken) IsExpired() bool {
	return time.Now().After(t.Expiry)
}

// IsValid 检查令牌是否有效
func (t *OAuthToken) IsValid() bool {
	return t.IsActive && !t.IsExpired()
}

// GetStatus 获取令牌状态
func (t *OAuthToken) GetStatus() OAuthTokenStatus {
	if !t.IsActive {
		return OAuthTokenStatusInactive
	}
	if t.IsExpired() {
		return OAuthTokenStatusExpired
	}
	return OAuthTokenStatusActive
}

// SetExpiry 根据expires_in设置过期时间
func (t *OAuthToken) SetExpiry() {
	if t.ExpiresIn > 0 {
		t.Expiry = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
	}
}
