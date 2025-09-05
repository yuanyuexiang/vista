# 微信授权API文档 - 前端驱动模式

## 概述

新架构采用前端驱动模式，后端只提供两个核心API接口：
1. **通过code获取用户信息** - `POST /api/wechat/auth`
2. **查询用户授权状态** - `GET /api/user/{openid}`

## API接口详情

### 1. 通过授权码获取用户信息

**接口地址：** `POST /api/wechat/auth`

**请求头：**
```
Content-Type: application/json
```

**请求参数：**
```json
{
  "code": "微信授权码",
  "state": "可选的状态参数"
}
```

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "微信登录成功",
    "user_info": {
      "openid": "o3yLZ14BG3RwIQivHu9qWhs4b6gg",
      "nickname": "tom",
      "headimgurl": "https://thirdwx.qlogo.cn/mmopen/...",
      "sex": 0,
      "language": "",
      "country": "",
      "province": "",
      "city": "",
      "privilege": [],
      "login_time": 1757062690
    }
  }
}
```

**错误响应：**
```json
{
  "code": 400,
  "message": "请求参数错误: Key: 'WechatAuthRequest.Code' Error:Field validation for 'Code' failed on the 'required' tag"
}
```

---

### 2. 查询用户授权状态

**接口地址：** `GET /api/user/{openid}`

**路径参数：**
- `openid`: 微信用户的OpenID

**用户存在且授权有效响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exists": true,
    "need_auth": false,
    "user_info": {
      "openid": "o3yLZ14BG3RwIQivHu9qWhs4b6gg",
      "nickname": "tom",
      "headimgurl": "https://thirdwx.qlogo.cn/mmopen/...",
      "sex": 0,
      "language": "",
      "country": "",
      "province": "",
      "city": "",
      "privilege": [],
      "login_time": 1757062690,
      "expires_at": 1757149090
    }
  }
}
```

**用户不存在响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exists": false,
    "need_auth": true,
    "message": "用户不存在，需要重新授权"
  }
}
```

**授权已过期响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exists": true,
    "need_auth": true,
    "user_info": { ... }
  }
}
```

---

## 前端集成指南

### 1. 构建微信授权链接

```javascript
function buildWechatAuthURL() {
    const appId = 'wx1eb05232cfbb49f7'; // 您的微信AppID
    const redirectURI = encodeURIComponent(window.location.href.split('?')[0]);
    const state = generateState();
    
    return `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appId}&redirect_uri=${redirectURI}&response_type=code&scope=snsapi_userinfo&state=${state}#wechat_redirect`;
}
```

### 2. 处理微信回调

```javascript
async function handleWechatCallback(code, state) {
    try {
        const response = await fetch('/api/wechat/auth', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                code: code,
                state: state
            })
        });

        const data = await response.json();
        
        if (data.code === 200) {
            const userInfo = data.data.user_info;
            // 保存用户信息到本地存储
            localStorage.setItem('wechat_user_info', JSON.stringify(userInfo));
            localStorage.setItem('wechat_openid', userInfo.openid);
        }
    } catch (error) {
        console.error('授权处理失败:', error);
    }
}
```

### 3. 检查授权状态

```javascript
async function checkAuthStatus(openid) {
    try {
        const response = await fetch(`/api/user/${openid}`);
        const data = await response.json();
        
        if (data.code === 200) {
            const result = data.data;
            
            if (!result.exists || result.need_auth) {
                // 需要重新授权
                startWechatAuth();
            } else {
                // 授权有效，显示用户信息
                displayUserInfo(result.user_info);
            }
        }
    } catch (error) {
        console.error('状态检查失败:', error);
    }
}
```

### 4. 完整流程示例

```javascript
// 页面加载时检查
window.onload = function() {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');
    
    if (code) {
        // 处理微信回调
        const state = urlParams.get('state');
        handleWechatCallback(code, state);
    } else {
        // 检查本地是否有openid
        const openid = localStorage.getItem('wechat_openid');
        if (openid) {
            checkAuthStatus(openid);
        } else {
            // 显示登录按钮
            showLoginButton();
        }
    }
};
```

---

## 授权状态说明

### 授权有效期
- 用户信息有效期：**24小时**
- 过期后需要重新授权
- `expires_at` 字段表示过期时间戳

### 状态字段含义
- `exists`: 用户在数据库中是否存在
- `need_auth`: 是否需要重新授权
  - `true`: 用户不存在或授权已过期
  - `false`: 授权有效，可直接使用

---

## 演示页面

### 前端驱动模式演示
**URL:** `https://carture.matrix-net.tech/wechat/frontend-driven`

**特性：**
- ✅ 完全由前端控制授权流程
- ✅ 智能检测授权状态
- ✅ 自动处理过期重新授权
- ✅ 支持强制重新授权

### 使用步骤
1. 在微信中打开演示页面
2. 点击"开始微信授权"按钮
3. 完成微信授权流程
4. 自动显示用户信息
5. 24小时后自动提示重新授权

---

## 错误处理

### 常见错误码
- `400`: 请求参数错误
- `500`: 服务器内部错误（微信API调用失败等）

### 前端错误处理建议
```javascript
try {
    const response = await fetch('/api/wechat/auth', { ... });
    const data = await response.json();
    
    if (data.code !== 200) {
        // 显示错误信息
        showError(data.message);
        return;
    }
    
    // 处理成功响应
    handleSuccess(data.data);
} catch (error) {
    // 网络错误处理
    showError('网络连接失败，请重试');
}
```

---

## 安全考虑

1. **State参数验证** - 防止CSRF攻击
2. **Code一次性使用** - 微信授权码只能使用一次
3. **授权过期机制** - 24小时自动过期，提高安全性
4. **HTTPS传输** - 确保数据传输安全

---

## 架构优势

### 🎯 前端控制
- 授权流程完全由前端控制
- 适配任何前端框架（React、Vue、Angular等）
- 易于自定义用户体验

### ⚡ 后端简洁
- 只需维护两个核心API接口
- 逻辑清晰，易于维护
- 降低后端复杂度

### 🔄 智能状态管理
- 自动检测授权过期
- 支持强制重新授权
- 本地存储优化

### 🚀 易于扩展
- 可轻松添加其他社交登录
- 支持多种认证方式
- 便于集成到现有系统
