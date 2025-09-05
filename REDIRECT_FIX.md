# 微信授权重定向修复说明

## 问题描述
用户反馈测试后发现微信授权完成后跳转到了根路径 `https://carture.matrix-net.tech/?wechat_auth=success&user_info=...` 而不是预期的 `/wechat/demo` 路径。

## 原因分析
问题在于微信授权流程中 `redirect_url` 参数没有正确传递：

1. 前端 JavaScript 正确发送了 `redirect_url` 参数
2. 但后端的 `WechatAuth` 处理器没有保存这个参数
3. 在微信回调时无法获取到原始的重定向URL

## 解决方案

### 1. 修改 `utils/wechat.go` 
添加了新的函数来处理包含重定向URL的state参数：

```go
// StateData 包含state参数中的数据
type StateData struct {
    Random      string `json:"random"`
    RedirectURL string `json:"redirect_url,omitempty"`
    Timestamp   int64  `json:"timestamp"`
}

// GenerateStateWithRedirect 生成包含重定向URL的state参数
func GenerateStateWithRedirect(redirectURL string) string

// ParseStateData 解析state参数中的数据
func ParseStateData(state string) (*StateData, error)
```

### 2. 修改 `WechatAuth` 处理器
```go
func WechatAuth(c *gin.Context) {
    // 获取前端传递的重定向URL
    redirectURL := c.Query("redirect_url")
    
    // 生成包含重定向URL的state参数
    state := utils.GenerateStateWithRedirect(redirectURL)
    
    // ... 其他代码
}
```

### 3. 修改 `WechatCallback` 处理器
```go
func WechatCallback(c *gin.Context) {
    // ... 其他代码
    
    // 解析state参数获取重定向URL
    stateData, err := utils.ParseStateData(state)
    if err != nil {
        fmt.Printf("Error parsing state data: %v\n", err)
    }

    // 使用state中的重定向URL
    redirectURL := ""
    if stateData != nil && stateData.RedirectURL != "" {
        redirectURL = stateData.RedirectURL
    }
    // 降级策略：查询参数 -> 默认值
    if redirectURL == "" {
        redirectURL = c.Query("redirect_url")
    }
    if redirectURL == "" {
        redirectURL = "https://carture.matrix-net.tech/"
    }
    
    // ... 其他代码
}
```

## 技术实现

### State参数编码
- 将重定向URL和随机字符串打包成JSON
- 使用Base64编码传递给微信
- 微信回调时解码获取原始重定向URL

### 降级策略
1. 优先使用state中的重定向URL
2. 如果解析失败，使用查询参数中的redirect_url
3. 如果都没有，使用默认的根路径

### 安全考虑
- state参数包含时间戳，可用于防止重放攻击
- 随机字符串防止CSRF攻击
- Base64编码确保URL安全传输

## 测试验证

### 测试用例
1. 从 `/wechat/demo` 发起授权，应该返回 `/wechat/demo`
2. 直接访问授权链接不带参数，应该返回根路径
3. 带有自定义redirect_url参数，应该返回指定路径

### 预期结果
```
授权前: https://carture.matrix-net.tech/wechat/demo
授权后: https://carture.matrix-net.tech/wechat/demo?wechat_auth=success&user_info=...
```

## 部署说明

修复后需要重新部署后端服务，前端 `wechat_demo.html` 无需修改，因为JavaScript已经正确发送了 `redirect_url` 参数。

## 文件变更列表

- ✅ `pkg/utils/wechat.go` - 新增state参数处理函数
- ✅ `internal/handler/wechat.go` - 修改授权和回调处理逻辑  
- ✅ `wechat_demo.html` - 已有正确的JavaScript实现
- ✅ `internal/router/router.go` - 路由配置保持不变
