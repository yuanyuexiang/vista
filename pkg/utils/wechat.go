package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateState 生成随机 state 参数
func GenerateState() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选
		return fmt.Sprintf("state_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

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
