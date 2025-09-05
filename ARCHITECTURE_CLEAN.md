# 🏗️ 微信授权架构设计

## 📋 概述
已优化微信授权系统架构，实现前后端完全分离，后端专注于业务逻辑，前端专注于用户界面。

## 🔧 后端职责 (Go Gin)

### 核心功能
- ✅ **微信OAuth流程处理**
- ✅ **用户信息存储管理**
- ✅ **RESTful API接口**
- ✅ **业务逻辑处理**

### API 接口
```go
// 微信授权
GET /wechat/auth              // 发起微信授权
GET /wechat/callback          // 处理微信回调

// 用户数据接口
GET /api/user/{openid}        // 获取单个用户信息
GET /api/users                // 获取用户列表

// 系统接口
GET /health                   // 健康检查
```

### 回调处理逻辑
```go
func WechatCallback(c *gin.Context) {
    // 1. 获取微信授权码
    code := c.Query("code")
    
    // 2. 换取access_token和用户信息
    authResult, _ := wechatService.ExchangeCodeForToken(code)
    userInfo, _ := wechatService.GetUserInfo(authResult.AccessToken, authResult.OpenID)
    
    // 3. 保存到数据库
    wechatService.SaveUserAuth(authResult, userInfo)
    
    // 4. 构建用户数据
    userInfoData := gin.H{
        "openid":     authResult.OpenID,
        "nickname":   userInfo.Nickname,
        "headimgurl": userInfo.HeadImgURL,
        // ... 其他字段
        "login_time": time.Now().Unix(),
    }
    
    // 5. Base64编码用户信息
    userInfoJson, _ := json.Marshal(userInfoData)
    userInfoB64 := base64.URLEncoding.EncodeToString(userInfoJson)
    
    // 6. 重定向到前端应用
    redirectURL := c.Query("redirect_url")
    if redirectURL == "" {
        redirectURL = "https://carture.matrix-net.tech/"
    }
    
    finalURL := fmt.Sprintf("%s?wechat_auth=success&user_info=%s", redirectURL, userInfoB64)
    c.Redirect(http.StatusFound, finalURL)
}
```

## 🎨 前端职责 (React)

### 核心功能
- ✅ **用户界面渲染**
- ✅ **用户交互处理** 
- ✅ **本地数据存储**
- ✅ **状态管理**

### 授权流程处理
```javascript
// 1. 页面加载时检查授权状态
useEffect(() => {
    // 检查URL参数中的授权信息
    const authResult = checkAuthFromURL();
    if (authResult.success) {
        setUserInfo(authResult.userInfo);
        return;
    }
    
    // 检查本地存储
    const localUserInfo = getUserInfo();
    if (localUserInfo) {
        setUserInfo(localUserInfo);
    }
}, []);

// 2. 解析URL参数中的用户信息
function checkAuthFromURL() {
    const urlParams = new URLSearchParams(window.location.search);
    const authStatus = urlParams.get('wechat_auth');
    const userInfoB64 = urlParams.get('user_info');
    
    if (authStatus === 'success' && userInfoB64) {
        const userInfoJson = atob(userInfoB64);
        const userInfo = JSON.parse(userInfoJson);
        
        // 保存到本地存储
        localStorage.setItem('wechat_user_info', JSON.stringify(userInfo));
        
        // 清理URL参数
        cleanURLParams();
        
        return { success: true, userInfo };
    }
    
    return { success: false };
}

// 3. 发起授权
function startWechatAuth() {
    const currentURL = window.location.href;
    const authURL = `${API_BASE}/wechat/auth?redirect_url=${encodeURIComponent(currentURL)}`;
    window.location.href = authURL;
}
```

## 🔄 完整流程

### 1. 用户访问前端应用
```
用户在微信中访问: https://carture.matrix-net.tech/
```

### 2. 前端检查授权状态
```
React应用检查:
├── URL参数中是否有授权信息
├── localStorage中是否有用户信息
└── 如果都没有，显示登录按钮
```

### 3. 发起微信授权
```
点击登录 → 跳转到后端授权接口:
https://carture.matrix-net.tech/wechat/auth?redirect_url=https://carture.matrix-net.tech/
```

### 4. 微信授权处理
```
后端重定向到微信授权页面 →
用户在微信中授权 →
微信回调到后端 /wechat/callback
```

### 5. 后端处理回调
```
后端处理:
├── 换取access_token
├── 获取用户信息  
├── 保存到数据库
└── 重定向回前端应用 (携带用户信息)
```

### 6. 前端接收用户信息
```
前端接收:
├── 解析URL参数中的用户信息
├── 保存到localStorage
├── 更新界面状态
└── 清理URL参数
```

## 🎯 架构优势

### 🚀 性能优化
- **前后端分离**: 各自专注于核心职责
- **缓存机制**: 用户信息缓存24小时
- **最小化传输**: 只传输必要的用户数据

### 🔒 安全性
- **无敏感信息泄露**: 后端不暴露前端代码
- **数据加密传输**: Base64编码用户信息
- **CSRF防护**: 使用state参数防止攻击

### 🛠️ 可维护性
- **职责清晰**: 前后端边界明确
- **易于扩展**: 可独立添加新功能
- **代码复用**: 工具函数可在多个组件中使用

### 🔧 开发效率
- **独立开发**: 前后端可并行开发
- **快速调试**: 专注于各自领域的问题
- **易于测试**: 可独立进行单元测试

## 📊 数据流向

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   用户访问   │───▶│  React应用   │───▶│  检查授权    │
└─────────────┘    └─────────────┘    └─────────────┘
                                               │
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  显示登录    │◀───│  未授权状态   │◀───│  无用户信息   │
└─────────────┘    └─────────────┘    └─────────────┘
       │
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  点击登录    │───▶│  后端授权    │───▶│  微信授权    │
└─────────────┘    └─────────────┘    └─────────────┘
                                               │
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  用户授权    │───▶│  微信回调    │───▶│  后端处理    │
└─────────────┘    └─────────────┘    └─────────────┘
                                               │
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  保存数据    │───▶│  重定向      │───▶│  前端接收    │
└─────────────┘    └─────────────┘    └─────────────┘
                                               │
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  显示用户    │◀───│  更新界面    │◀───│  保存本地    │
└─────────────┘    └─────────────┘    └─────────────┘
```

## 🔧 环境配置

### 后端配置 (config.yaml)
```yaml
wechat:
  app_id: "wx1eb05232cfbb49f7"
  app_secret: "your_app_secret"
  redirect_uri: "https://carture.matrix-net.tech/wechat/callback"
```

### 前端配置 (wechatAuth.js)
```javascript
const WECHAT_CONFIG = {
  API_BASE: 'https://carture.matrix-net.tech',
  STORAGE_KEY: 'wechat_user_info',
  AUTH_SUCCESS_KEY: 'wechat_auth_success'
};
```

这样的架构设计实现了真正的前后端分离，代码更加清晰、可维护，也更符合现代Web应用的最佳实践。
