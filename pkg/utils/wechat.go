package utils

import (
	"time"
)

// IsValidOpenID 验证 OpenID 格式
func IsValidOpenID(openid string) bool {
	if len(openid) < 10 || len(openid) > 128 {
		return false
	}
	// 可以添加更多的格式验证
	return true
}

// IsTokenExpired 检查令牌是否过期
func IsTokenExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}
