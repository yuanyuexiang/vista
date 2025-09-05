# 微信自动授权演示使用说明

## 功能概述

这个演示页面实现了智能的微信登录检测和自动授权功能，当用户在微信中访问 `https://carture.matrix-net.tech/wechat/demo` 时会自动处理登录流程。

## 自动处理流程

### 1. 智能检测
页面加载时会按以下顺序检测：

1. **URL回调检测** - 检查是否是微信授权回调
2. **本地缓存检测** - 检查是否有有效的登录信息（24小时有效期）
3. **微信环境检测** - 验证是否在微信浏览器中
4. **自动授权跳转** - 未登录时自动跳转到微信授权页面

### 2. 用户体验
- ✅ **已登录用户**：直接显示用户信息，无需重复授权
- 🔄 **未登录用户**：在微信中自动跳转授权，非微信环境显示提示
- ⏰ **过期处理**：登录信息过期后自动清除并重新授权
- 🔄 **URL清理**：授权完成后自动清理URL参数，保持页面整洁

## 技术特性

### 智能缓存机制
```javascript
// 用户信息自动保存到 localStorage
localStorage.setItem('wechat_user_info', JSON.stringify(userInfo));
localStorage.setItem('wechat_auth_success', 'true');

// 24小时过期检查
const loginTime = userInfo.login_time || 0;
const now = Math.floor(Date.now() / 1000);
if (now - loginTime < 24 * 60 * 60) {
    // 信息有效
} else {
    // 自动清除过期信息
}
```

### 微信环境检测
```javascript
function isWechatBrowser() {
    const ua = navigator.userAgent.toLowerCase();
    return /micromessenger/.test(ua);
}
```

### Base64数据传输
授权成功后通过URL参数传输Base64编码的用户信息：
```
https://carture.matrix-net.tech/wechat/demo?wechat_auth=success&user_info=eyJvcGVuaWQiOi...
```

## 使用方法

### 1. 微信中直接访问
```
https://carture.matrix-net.tech/wechat/demo
```

### 2. 测试流程
1. 首次访问：自动跳转微信授权
2. 授权完成：自动返回显示用户信息  
3. 再次访问：直接显示已保存的信息
4. 清除信息：点击"清除用户信息"按钮
5. 重新授权：点击"重新授权"按钮

### 3. 非微信环境
在非微信浏览器中访问会显示提示信息和手动登录选项。

## API 接口

系统提供RESTful API供前端调用：

```bash
# 获取单个用户信息
GET /api/user/{openid}

# 获取所有用户列表
GET /api/users

# 微信授权流程
GET /wechat/auth           # 发起授权
GET /wechat/callback       # 授权回调
```

## 状态提示

页面会实时显示不同颜色的状态信息：
- 🔵 **蓝色** - 信息提示（检查中、跳转中）
- 🟢 **绿色** - 成功状态（登录成功、操作完成）
- 🔴 **红色** - 错误提示（需要微信环境、授权失败）

## 开发集成

### React项目集成
```javascript
// 检查微信登录状态
const checkWechatAuth = () => {
    const userInfo = localStorage.getItem('wechat_user_info');
    if (userInfo) {
        const user = JSON.parse(userInfo);
        // 使用用户信息
        return user;
    }
    return null;
};

// 发起微信授权
const startWechatAuth = () => {
    const currentURL = window.location.href;
    const authURL = `https://carture.matrix-net.tech/wechat/auth?redirect_url=${encodeURIComponent(currentURL)}`;
    window.location.href = authURL;
};
```

### Vue项目集成
```javascript
// 在 mounted 钩子中检查登录状态
mounted() {
    this.checkWechatLogin();
},

methods: {
    checkWechatLogin() {
        const userInfo = localStorage.getItem('wechat_user_info');
        if (userInfo) {
            this.user = JSON.parse(userInfo);
            this.isLoggedIn = true;
        }
    }
}
```

## 配置说明

### 环境配置
- **开发环境**：使用微信测试号 (wx1eb05232cfbb49f7)
- **生产环境**：需要配置正式的微信公众号
- **域名配置**：carture.matrix-net.tech（需在微信后台配置）

### 安全考虑
- ✅ 用户信息通过HTTPS传输
- ✅ 支持过期时间检查
- ✅ URL参数自动清理
- ✅ 错误处理和fallback机制

## 故障排除

### 常见问题
1. **授权失败** - 检查微信测试号配置和域名设置
2. **信息丢失** - 检查localStorage是否被清除
3. **跳转异常** - 确认redirect_url参数正确
4. **非微信环境** - 使用手动登录选项或在微信中打开

### 调试工具
- 浏览器开发者工具查看localStorage
- 网络面板检查API请求
- 控制台查看错误信息
- 页面状态提示实时反馈
