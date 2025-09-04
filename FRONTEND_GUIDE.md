# Vista 微信授权服务 - 前端集成指南

## 📋 概述

基于 REST API 的微信公众号授权服务，提供微信 OAuth2 授权登录功能。

## 🚀 服务地址

- **微信授权接口**: `http://localhost:8080/wechat/auth` (发起授权跳转)
- **微信回调接口**: `http://localhost:8080/wechat/callback` (微信授权回调)
- **健康检查**: `http://localhost:8080/health` (服务状态检查)

## 🔄 完整授权流程

### 流程说明

1. **前端引导用户跳转**: 用户点击"微信登录"按钮，跳转到后端授权接口
2. **后端重定向到微信**: 后端构建微信授权 URL，重定向用户到微信授权页面
3. **用户在微信授权**: 用户在微信页面完成授权
4. **微信回调到后端**: 微信将 code 回调到后端 callback 接口
5. **后端处理授权**: 后端用 code 换取 access_token，获取用户信息
6. **重定向回前端**: 后端重定向回前端，带上用户信息

### 1. 前端发起授权

直接跳转到后端授权接口即可：

```javascript
function startWechatAuth() {
  // 直接跳转到后端授权接口
  window.location.href = 'http://localhost:8080/wechat/auth';
}
```

```html
<!-- 微信登录按钮 -->
<button onclick="startWechatAuth()">微信登录</button>

<!-- 或者直接使用链接 -->
<a href="http://localhost:8080/wechat/auth">微信登录</a>
```

### 2. 处理授权成功回调

后端会重定向到前端页面，带上用户信息：

```javascript
// 在授权成功页面（如 /auth/success）处理回调
function handleAuthSuccess() {
  const params = new URLSearchParams(window.location.search);
  const openid = params.get('openid');
  const nickname = params.get('nickname');
  
  if (openid) {
    // 保存用户信息
    const userInfo = {
      openid: openid,
      nickname: decodeURIComponent(nickname || ''),
      loginTime: new Date().toISOString()
    };
    
    localStorage.setItem('wechat_user', JSON.stringify(userInfo));
    
    // 跳转到主页面
    window.location.href = '/dashboard';
  } else {
    // 授权失败处理
    alert('微信授权失败，请重试');
    window.location.href = '/login';
  }
}

// 页面加载时自动处理
window.onload = function() {
  if (window.location.pathname === '/auth/success') {
    handleAuthSuccess();
  }
};
```

### 3. 检查登录状态

```javascript
// 检查用户是否已登录
function isLoggedIn() {
  const userInfo = localStorage.getItem('wechat_user');
  return userInfo !== null;
}

// 获取当前用户信息
function getCurrentUser() {
  const userInfo = localStorage.getItem('wechat_user');
  return userInfo ? JSON.parse(userInfo) : null;
}

// 退出登录
function logout() {
  localStorage.removeItem('wechat_user');
  window.location.href = '/login';
}
```

### 4. API 调用示例（可选）

如果需要调用其他 API，可以携带用户信息：

```javascript
// 调用受保护的 API
async function callProtectedAPI(endpoint, data = {}) {
  const user = getCurrentUser();
  if (!user) {
    window.location.href = '/login';
    return;
  }
  
  try {
    const response = await fetch(`http://localhost:8080${endpoint}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-OpenID': user.openid, // 携带用户标识
      },
      body: JSON.stringify(data)
    });
    
    if (!response.ok) {
      throw new Error('API call failed');
    }
    
    return await response.json();
  } catch (error) {
    console.error('API 调用失败:', error);
    throw error;
  }
}
```

## 🛠️ 实用工具类

```javascript
// 微信授权工具类
class WechatAuthManager {
  constructor(config = {}) {
    this.baseURL = config.baseURL || 'http://localhost:8080';
    this.frontendDomain = config.frontendDomain || 'http://localhost:3000';
  }
  
  // 发起微信授权
  startAuth() {
    window.location.href = `${this.baseURL}/wechat/auth`;
  }
  
  // 处理授权回调
  handleCallback() {
    const params = new URLSearchParams(window.location.search);
    const openid = params.get('openid');
    const nickname = params.get('nickname');
    
    if (!openid) {
      throw new Error('授权失败：未获取到用户信息');
    }
    
    const userInfo = {
      openid: openid,
      nickname: decodeURIComponent(nickname || ''),
      loginTime: new Date().toISOString()
    };
    
    this.saveUser(userInfo);
    return userInfo;
  }
  
  // 保存用户信息
  saveUser(userInfo) {
    localStorage.setItem('wechat_user', JSON.stringify(userInfo));
  }
  
  // 获取用户信息
  getUser() {
    const userInfo = localStorage.getItem('wechat_user');
    return userInfo ? JSON.parse(userInfo) : null;
  }
  
  // 检查登录状态
  isLoggedIn() {
    return this.getUser() !== null;
  }
  
  // 退出登录
  logout() {
    localStorage.removeItem('wechat_user');
  }
}

// 使用示例
const authManager = new WechatAuthManager({
  baseURL: 'http://localhost:8080',
  frontendDomain: 'http://localhost:3000'
});

// 发起授权
document.getElementById('wechat-login').onclick = () => {
  authManager.startAuth();
};

// 在回调页面处理授权结果
if (window.location.pathname === '/auth/success') {
  try {
    const user = authManager.handleCallback();
    console.log('登录成功:', user);
    window.location.href = '/dashboard';
  } catch (error) {
    console.error('登录失败:', error);
    window.location.href = '/login';
  }
}
```

## 🔧 测试

### 1. 启动后端服务

```bash
cd vista
go run main.go
```

服务启动后显示：
```
Server starting on port 8080
WeChat Auth URL: http://localhost:8080/wechat/auth
WeChat Callback URL: http://localhost:8080/wechat/callback
```

### 2. 测试授权流程

1. 浏览器访问：`http://localhost:8080/wechat/auth`
2. 应该会重定向到微信授权页面
3. 授权完成后会回调到 `/wechat/callback`
4. 最终重定向到前端页面

### 3. 健康检查

```bash
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok",
  "service": "vista-wechat-auth"
}
```

## ⚠️ 注意事项

### 配置要求
- 确保 `config.yaml` 中配置了正确的微信 AppID 和 AppSecret
- RedirectURI 必须与微信公众号后台配置一致
- 生产环境使用 HTTPS

### 安全考虑
- state 参数用于防止 CSRF 攻击
- 不要在前端存储敏感的 access_token
- 生产环境应该生成 JWT token 而不是直接传递微信信息

### 错误处理
- 处理授权被拒绝的情况
- 处理网络错误和超时
- 提供用户友好的错误提示

### 兼容性
- 确保在微信浏览器中正常工作
- 测试不同版本的微信客户端
- 移动端和桌面端兼容性测试

## 📱 前端页面示例

### 登录页面 (login.html)
```html
<!DOCTYPE html>
<html>
<head>
    <title>微信登录</title>
</head>
<body>
    <div style="text-align: center; margin-top: 100px;">
        <h2>欢迎使用 Vista 服务</h2>
        <a href="http://localhost:8080/wechat/auth" 
           style="display: inline-block; padding: 10px 20px; background: #07c160; color: white; text-decoration: none; border-radius: 5px;">
            微信登录
        </a>
    </div>
</body>
</html>
```

### 授权成功页面 (auth/success.html)
```html
<!DOCTYPE html>
<html>
<head>
    <title>授权成功</title>
</head>
<body>
    <div style="text-align: center; margin-top: 100px;">
        <h2>正在处理授权信息...</h2>
        <p>请稍候，即将跳转到主页</p>
    </div>
    
    <script>
        // 自动处理授权回调
        const params = new URLSearchParams(window.location.search);
        const openid = params.get('openid');
        const nickname = params.get('nickname');
        
        if (openid) {
            const userInfo = {
                openid: openid,
                nickname: decodeURIComponent(nickname || ''),
                loginTime: new Date().toISOString()
            };
            
            localStorage.setItem('wechat_user', JSON.stringify(userInfo));
            
            setTimeout(() => {
                window.location.href = '/dashboard.html';
            }, 2000);
        } else {
            alert('授权失败，请重试');
            window.location.href = '/login.html';
        }
    </script>
</body>
</html>
```
