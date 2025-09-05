package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// StateData 包含state参数中的数据
type StateData struct {
	Random      string `json:"random"`
	RedirectURL string `json:"redirect_url,omitempty"`
	Timestamp   int64  `json:"timestamp"`
}

// GenerateState 生成随机 state 参数
func GenerateState() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选
		return fmt.Sprintf("state_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// GenerateStateWithRedirect 生成包含重定向URL的state参数
func GenerateStateWithRedirect(redirectURL string) string {
	// 生成随机字符串
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return GenerateState() // 降级到简单state
	}

	// 构建state数据
	stateData := StateData{
		Random:      hex.EncodeToString(bytes),
		RedirectURL: redirectURL,
		Timestamp:   time.Now().Unix(),
	}

	// 序列化为JSON并编码为base64
	jsonData, err := json.Marshal(stateData)
	if err != nil {
		return GenerateState() // 降级到简单state
	}

	return base64.URLEncoding.EncodeToString(jsonData)
}

// ParseStateData 解析state参数中的数据
func ParseStateData(state string) (*StateData, error) {
	// 尝试base64解码
	jsonData, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		// 如果解码失败，可能是简单的state格式
		return &StateData{
			Random:    state,
			Timestamp: time.Now().Unix(),
		}, nil
	}

	// 解析JSON
	var stateData StateData
	if err := json.Unmarshal(jsonData, &stateData); err != nil {
		return &StateData{
			Random:    state,
			Timestamp: time.Now().Unix(),
		}, nil
	}

	return &stateData, nil
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
