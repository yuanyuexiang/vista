# 前端获取微信用户信息指南

## 概述
当微信浏览器发起回调请求 `https://carture.matrix-net.tech/wechat/callback?code=xxx&state=xxx` 时，前端有多种方式可以获取到用户信息。

## 方法1: HTML页面 + localStorage (推荐)

### 后端实现
回调接口返回包含用户信息的HTML页面，页面中的JavaScript会：
- 将用户信息保存到 `localStorage`
- 通过 `postMessage` 发送给父页面 
- 在微信浏览器中自动关闭窗口

### 前端获取
```javascript
// 方式1: 从 localStorage 读取
const userInfo = localStorage.getItem('wechat_user_info');
if (userInfo) {
    const user = JSON.parse(userInfo);
    console.log('用户信息:', user);
}

// 方式2: 监听 postMessage 事件
window.addEventListener('message', function(event) {
    if (event.data.type === 'WECHAT_LOGIN_SUCCESS') {
        const userInfo = event.data.data;
        console.log('接收到用户信息:', userInfo);
    }
});
```

## 方法2: API接口获取

### 接口列表
```
GET /api/user/{openid}  - 获取单个用户信息
GET /api/users          - 获取所有用户列表
```

### 前端调用
```javascript
// 获取特定用户信息
async function getUserInfo(openid) {
    try {
        const response = await fetch(`/api/user/${openid}`);
        const data = await response.json();
        if (data.code === 200) {
            return data.data.user_info;
        }
    } catch (error) {
        console.error('获取用户信息失败:', error);
    }
}

// 获取所有用户列表
async function getAllUsers() {
    try {
        const response = await fetch('/api/users');
        const data = await response.json();
        if (data.code === 200) {
            return data.data.users;
        }
    } catch (error) {
        console.error('获取用户列表失败:', error);
    }
}
```

## 方法3: 弹窗登录

### 实现方式
```javascript
function openLoginWindow() {
    const popup = window.open(
        'https://carture.matrix-net.tech/wechat/auth',
        'wechat_login',
        'width=400,height=600,scrollbars=yes,resizable=yes'
    );

    // 检查弹窗是否被关闭
    const checkClosed = setInterval(() => {
        if (popup.closed) {
            clearInterval(checkClosed);
            // 弹窗关闭后检查localStorage中的用户信息
            checkUserInfo();
        }
    }, 1000);
}
```

## 方法4: iframe嵌入

### HTML结构
```html
<iframe id="loginIframe" src="https://carture.matrix-net.tech/wechat/auth" 
        width="100%" height="500" frameborder="0"></iframe>
```

### JavaScript监听
```javascript
window.addEventListener('message', function(event) {
    if (event.data.type === 'WECHAT_LOGIN_SUCCESS') {
        const userInfo = event.data.data;
        // 处理用户信息
        handleUserInfo(userInfo);
        
        // 关闭iframe
        document.getElementById('loginIframe').src = '';
    }
});
```

## 数据结构

### 用户信息格式
```json
{
    "openid": "o3yLZ14BG3RwIQivHu9qWhs4b6gg",
    "nickname": "tom",
    "sex": 0,
    "language": "",
    "city": "",
    "province": "",
    "country": "",
    "headimgurl": "https://thirdwx.qlogo.cn/mmopen/.../132",
    "privilege": []
}
```

### 性别字段说明
- 0: 未知
- 1: 男
- 2: 女

## 测试页面

### 本地测试
访问: `http://localhost:8080/demo`

### 线上测试  
访问: `https://carture.matrix-net.tech/demo`

## 注意事项

1. **跨域问题**: 确保正确设置CORS头
2. **HTTPS要求**: 微信要求回调地址必须是HTTPS
3. **域名白名单**: 需要在微信公众平台配置授权域名
4. **localStorage生命周期**: 数据会持久保存，注意及时清理
5. **postMessage安全**: 验证消息来源，防止XSS攻击

## 推荐方案

对于大部分应用场景，推荐使用 **方法1 (HTML页面 + localStorage)**：

**优势:**
- ✅ 简单易用，无需复杂的通信机制
- ✅ 数据持久化，页面刷新不丢失  
- ✅ 支持多种登录方式（直接跳转、弹窗、iframe）
- ✅ 在微信浏览器中体验良好

**使用步骤:**
1. 引导用户点击登录按钮
2. 跳转到微信授权页面
3. 授权完成后显示成功页面
4. 前端从localStorage读取用户信息

这种方案在微信生态中是最常见和稳定的实现方式。
