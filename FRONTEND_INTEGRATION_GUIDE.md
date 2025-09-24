# 📱 微信OAuth前端集成指南

## 🎯 概述

本文档为前端开发者提供完整的微信OAuth集成方案。采用**前端驱动模式**，后端只提供2个核心API，前端完全控制授权流程。

## 🏗️ 架构设计

```
前端应用 ──────────────────── 微信OAuth服务器
    │                           │
    │ 1. 构建授权链接              │
    │ ────────────────────────► │
    │                           │
    │ 2. 用户授权后获取code        │
    │ ◄──────────────────────── │
    │                           │
    │ 3. 发送code到后端           │
    ▼                           │
Vista后端API                    │
    │                           │
    │ 4. 通过code获取用户信息      │
    │ ────────────────────────► │
    │ ◄──────────────────────── │
    │                           │
    │ 5. 返回用户信息给前端        │
    ▼                           │
前端应用                         │
```

## 🔗 API接口

### 基础配置

```javascript
const CONFIG = {
    API_BASE: 'https://carture.kcbaotech.com',
    WECHAT_APP_ID: 'wx1eb05232cfbb49f7', // 微信AppID
};
```

### 1. 授权接口

**POST** `/vista/wechat/api/auth`

通过微信授权码获取用户信息并保存到数据库。

**请求头：**
```
Content-Type: application/json
```

**请求体：**
```json
{
  "code": "微信授权码",
  "state": "可选的状态参数"
}
```

**响应：**
```json
{
  "code": 200,
  "message": "授权成功",
  "data": {
    "user_info": {
      "openid": "用户唯一标识",
      "nickname": "用户昵称", 
      "sex": 1,
      "province": "省份",
      "city": "城市",
      "country": "国家",
      "headimgurl": "头像URL",
      "privilege": [],
      "unionid": "UnionID",
      "expires_at": 1725552000,
      "created_at": "2024-09-05T10:00:00Z",
      "updated_at": "2024-09-05T10:00:00Z"
    }
  }
}
```

### 2. 用户信息查询接口

**GET** `/vista/wechat/api/user/{openid}`

查询用户信息和授权状态，智能判断是否需要重新授权。

**响应：**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "exists": true,
    "need_auth": false,
    "user_info": {
      // 用户详细信息（同上）
    }
  }
}
```

**状态说明：**
- `exists: false` - 用户不存在，需要首次授权
- `need_auth: true` - 授权已过期（超过24小时），需要重新授权
- `need_auth: false` - 授权有效，可直接使用

## 🚀 前端集成步骤

### 1. 检测微信浏览器环境

```javascript
function isWechatBrowser() {
    const ua = navigator.userAgent.toLowerCase();
    return /micromessenger/.test(ua);
}
```

### 2. 构建微信授权URL

```javascript
function buildWechatAuthURL() {
    const redirectURI = encodeURIComponent(window.location.href.split('?')[0]);
    const state = 'state_' + Math.random().toString(36).substr(2, 16) + Date.now();
    
    return `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${CONFIG.WECHAT_APP_ID}&redirect_uri=${redirectURI}&response_type=code&scope=snsapi_userinfo&state=${state}#wechat_redirect`;
}
```

### 3. 启动授权流程

```javascript
function startWechatAuth() {
    if (!isWechatBrowser()) {
        alert('请在微信中打开此页面进行授权');
        return;
    }

    const authURL = buildWechatAuthURL();
    window.location.href = authURL;
}
```

### 4. 处理授权回调

```javascript
async function handleWechatCallback(code, state) {
    try {
        const response = await fetch(`${CONFIG.API_BASE}/vista/wechat/api/auth`, {
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
            
            // 清理URL参数
            const url = new URL(window.location);
            url.searchParams.delete('code');
            url.searchParams.delete('state');
            window.history.replaceState({}, document.title, url.toString());
            
            return userInfo;
        } else {
            throw new Error(data.message);
        }
    } catch (error) {
        console.error('授权处理失败:', error);
        throw error;
    }
}
```

### 5. 检查授权状态

```javascript
async function checkAuthStatus() {
    const openid = localStorage.getItem('wechat_openid');
    
    if (!openid) {
        return { needAuth: true, reason: '未登录' };
    }

    try {
        const response = await fetch(`${CONFIG.API_BASE}/vista/wechat/api/user/${openid}`);
        const data = await response.json();
        
        if (data.code === 200) {
            const result = data.data;
            
            if (!result.exists) {
                // 清除无效的本地数据
                localStorage.removeItem('wechat_user_info');
                localStorage.removeItem('wechat_openid');
                return { needAuth: true, reason: '用户不存在' };
            } 
            
            if (result.need_auth) {
                return { needAuth: true, reason: '授权已过期' };
            }
            
            return { 
                needAuth: false, 
                userInfo: result.user_info 
            };
        } else {
            return { needAuth: true, reason: '检查失败' };
        }
    } catch (error) {
        console.error('状态检查失败:', error);
        return { needAuth: true, reason: '网络错误' };
    }
}
```

## 📋 完整集成示例

### HTML页面结构

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>微信授权示例</title>
</head>
<body>
    <div id="app">
        <!-- 未登录状态 -->
        <div id="login-section">
            <button onclick="startWechatAuth()">微信登录</button>
        </div>
        
        <!-- 已登录状态 -->
        <div id="user-section" style="display: none;">
            <div id="user-info"></div>
            <button onclick="logout()">退出登录</button>
        </div>
    </div>

    <script>
        // 配置
        const CONFIG = {
            API_BASE: 'https://carture.kcbaotech.com',
            WECHAT_APP_ID: 'wx1eb05232cfbb49f7'
        };

        // 页面加载时检查授权状态
        window.onload = async function() {
            const urlParams = new URLSearchParams(window.location.search);
            const code = urlParams.get('code');
            const state = urlParams.get('state');
            
            if (code) {
                // 处理微信回调
                try {
                    const userInfo = await handleWechatCallback(code, state);
                    showUserInfo(userInfo);
                } catch (error) {
                    alert('授权失败: ' + error.message);
                }
            } else {
                // 检查现有授权状态
                const status = await checkAuthStatus();
                if (!status.needAuth) {
                    showUserInfo(status.userInfo);
                }
            }
        };

        // 显示用户信息
        function showUserInfo(userInfo) {
            document.getElementById('login-section').style.display = 'none';
            document.getElementById('user-section').style.display = 'block';
            document.getElementById('user-info').innerHTML = `
                <h3>欢迎, ${userInfo.nickname}</h3>
                <img src="${userInfo.headimgurl}" width="60" height="60" style="border-radius: 50%;">
                <p>OpenID: ${userInfo.openid}</p>
            `;
        }

        // 退出登录
        function logout() {
            localStorage.removeItem('wechat_user_info');
            localStorage.removeItem('wechat_openid');
            location.reload();
        }

        // 这里添加前面的工具函数...
    </script>
</body>
</html>
```

## 🎨 React集成示例

### Hook实现

```javascript
import { useState, useEffect } from 'react';

const useWechatAuth = () => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [needAuth, setNeedAuth] = useState(false);

    const CONFIG = {
        API_BASE: 'https://carture.kcbaotech.com',
        WECHAT_APP_ID: 'wx1eb05232cfbb49f7'
    };

    useEffect(() => {
        initAuth();
    }, []);

    const initAuth = async () => {
        const urlParams = new URLSearchParams(window.location.search);
        const code = urlParams.get('code');
        
        if (code) {
            // 处理微信回调
            await handleCallback(code);
        } else {
            // 检查现有授权
            await checkStatus();
        }
        
        setLoading(false);
    };

    const handleCallback = async (code) => {
        try {
            const response = await fetch(`${CONFIG.API_BASE}/vista/wechat/api/auth`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ code })
            });

            const data = await response.json();
            if (data.code === 200) {
                const userInfo = data.data.user_info;
                localStorage.setItem('wechat_user_info', JSON.stringify(userInfo));
                localStorage.setItem('wechat_openid', userInfo.openid);
                setUser(userInfo);
                setNeedAuth(false);
                
                // 清理URL
                window.history.replaceState({}, '', window.location.pathname);
            }
        } catch (error) {
            console.error('授权失败:', error);
            setNeedAuth(true);
        }
    };

    const checkStatus = async () => {
        const openid = localStorage.getItem('wechat_openid');
        if (!openid) {
            setNeedAuth(true);
            return;
        }

        try {
            const response = await fetch(`${CONFIG.API_BASE}/vista/wechat/api/user/${openid}`);
            const data = await response.json();
            
            if (data.code === 200 && data.data.exists && !data.data.need_auth) {
                setUser(data.data.user_info);
                setNeedAuth(false);
            } else {
                setNeedAuth(true);
            }
        } catch (error) {
            setNeedAuth(true);
        }
    };

    const startAuth = () => {
        const redirectURI = encodeURIComponent(window.location.href.split('?')[0]);
        const state = Date.now().toString();
        const authURL = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${CONFIG.WECHAT_APP_ID}&redirect_uri=${redirectURI}&response_type=code&scope=snsapi_userinfo&state=${state}#wechat_redirect`;
        window.location.href = authURL;
    };

    const logout = () => {
        localStorage.removeItem('wechat_user_info');
        localStorage.removeItem('wechat_openid');
        setUser(null);
        setNeedAuth(true);
    };

    return {
        user,
        loading,
        needAuth,
        startAuth,
        logout,
        checkStatus
    };
};

export default useWechatAuth;
```

### 组件使用

```javascript
import useWechatAuth from './useWechatAuth';

const App = () => {
    const { user, loading, needAuth, startAuth, logout } = useWechatAuth();

    if (loading) {
        return <div>加载中...</div>;
    }

    if (needAuth) {
        return (
            <div>
                <h2>请先登录</h2>
                <button onClick={startAuth}>微信登录</button>
            </div>
        );
    }

    return (
        <div>
            <h2>欢迎, {user.nickname}</h2>
            <img src={user.headimgurl} width="60" height="60" />
            <button onClick={logout}>退出登录</button>
        </div>
    );
};
```

## ⚠️ 注意事项

### 1. 安全考虑
- **不要在前端存储敏感信息**：只在localStorage中存储OpenID和基本用户信息
- **HTTPS必须**：微信OAuth要求所有回调地址必须使用HTTPS
- **域名白名单**：确保回调域名已在微信公众平台配置

### 2. 兼容性
- **微信浏览器专用**：OAuth流程只能在微信内置浏览器中进行
- **移动端优化**：确保页面在移动设备上正常显示
- **网络处理**：做好网络异常和超时处理

### 3. 用户体验
- **授权过期处理**：24小时后自动提示重新授权
- **加载状态**：在授权过程中显示加载提示
- **错误处理**：友好的错误提示和重试机制

### 4. 性能优化
- **本地缓存**：合理使用localStorage减少API调用
- **状态管理**：避免重复的授权状态检查
- **懒加载**：按需加载用户信息

## 🔧 测试与调试

### 本地测试
```bash
# 1. 启动服务
./vista

# 2. 访问测试页面
https://carture.kcbaotech.com/vista/wechat/frontend-driven
```

### 调试技巧
- 使用浏览器开发者工具查看网络请求
- 检查localStorage中的用户信息
- 查看console.log输出的调试信息

## 📚 相关资源

- [微信公众平台开发文档](https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html)
- [完整示例页面](https://carture.kcbaotech.com/vista/wechat/frontend-driven)
- [API接口文档](./API_DOCUMENTATION.md)

---

## 📞 技术支持

如有问题，请查看：
1. [常见问题解答](./FAQ.md) 
2. [API文档](./API_DOCUMENTATION.md)
3. 或提交Issue到项目仓库

---

*最后更新时间：2024年9月5日*
