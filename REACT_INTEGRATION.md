# React Expo 微信授权集成指南

## 🎯 目标
在React Expo应用中实现微信授权登录，用户在微信中访问 `https://carture.matrix-net.tech/` 时能够自动弹出授权，用户同意后获取用户信息并保存在本地。

## 📦 文件结构
```
your-react-app/
├── utils/
│   └── wechatAuth.js          # 微信授权工具函数
├── hooks/
│   └── useWechatAuth.js       # 微信授权Hook
├── components/
│   └── WechatLogin.jsx        # 微信登录组件
└── App.jsx                    # 主应用组件
```

## 🚀 快速开始

### 1. 安装依赖
```bash
# 如果你的项目还没有这些依赖，需要安装
npm install react react-native
# 或者
yarn add react react-native
```

### 2. 复制文件
将提供的以下文件复制到你的项目中：
- `utils/wechatAuth.js`
- `hooks/useWechatAuth.js` 
- `components/WechatLogin.jsx`

### 3. 配置后端地址
在 `utils/wechatAuth.js` 中确认配置：
```javascript
const WECHAT_CONFIG = {
  API_BASE: 'https://carture.matrix-net.tech', // 你的后端API地址
  STORAGE_KEY: 'wechat_user_info',
  AUTH_SUCCESS_KEY: 'wechat_auth_success'
};
```

### 4. 在应用中使用

#### 方式1: 使用WechatLogin组件 (推荐)
```jsx
import React from 'react';
import WechatLogin from './components/WechatLogin';

const App = () => {
  return (
    <div>
      <WechatLogin />
    </div>
  );
};

export default App;
```

#### 方式2: 使用useWechatAuth Hook
```jsx
import React from 'react';
import useWechatAuth from './hooks/useWechatAuth';

const MyComponent = () => {
  const {
    userInfo,
    isAuthenticated,
    isLoading,
    login,
    logout
  } = useWechatAuth({
    autoCheck: true,
    onAuthSuccess: (userInfo) => {
      console.log('登录成功:', userInfo);
    }
  });

  if (isLoading) return <div>加载中...</div>;

  if (!isAuthenticated) {
    return (
      <div>
        <button onClick={login}>微信登录</button>
      </div>
    );
  }

  return (
    <div>
      <h2>欢迎, {userInfo.nickname}!</h2>
      <button onClick={logout}>退出登录</button>
    </div>
  );
};
```

## 🔧 工作流程

### 1. 用户访问页面
用户在微信中访问 `https://carture.matrix-net.tech/`

### 2. 检查授权状态
React应用启动时自动检查：
- URL参数中是否有授权回调信息
- localStorage中是否有有效的用户信息

### 3. 发起授权
如果没有有效登录信息，显示登录按钮
```javascript
// 点击登录按钮时
startWechatAuth(); // 跳转到 /wechat/auth
```

### 4. 微信授权
用户在微信授权页面同意授权

### 5. 接收回调
微信回调到后端，后端处理后重定向回React应用：
```
https://carture.matrix-net.tech/?wechat_auth=success&user_info=base64编码的用户信息
```

### 6. 处理授权结果
React应用检测URL参数，解码用户信息并保存到localStorage

### 7. 更新界面
显示用户信息和登录后的内容

## 📱 在Expo中的特殊处理

### 1. 修改app.json
```json
{
  "expo": {
    "scheme": "your-app-scheme",
    "web": {
      "bundler": "metro"
    }
  }
}
```

### 2. 处理React Native样式
如果你使用的是React Native (Expo)，需要将CSS样式转换为StyleSheet：

```jsx
// 替换 CSS 样式为 React Native StyleSheet
import { StyleSheet } from 'react-native';

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
    backgroundColor: '#f5f5f5'
  },
  // ... 其他样式
});
```

### 3. 替换HTML元素
```jsx
// 替换 div -> View
// 替换 img -> Image  
// 替换 button -> TouchableOpacity + Text
// 替换 p -> Text

import { View, Text, Image, TouchableOpacity } from 'react-native';
```

## 🔐 数据安全

### 1. 用户信息加密
用户敏感信息在传输过程中使用base64编码，你可以根据需要增加加密：

```javascript
// 在 wechatAuth.js 中增加加密/解密函数
const encrypt = (text, key) => {
  // 使用你喜欢的加密算法
  return encryptedText;
};

const decrypt = (encryptedText, key) => {
  // 对应的解密算法
  return text;
};
```

### 2. Token刷新
用户信息设置24小时过期，过期后自动清除：

```javascript
// 检查信息是否过期（24小时）
const loginTime = userInfo.login_time || 0;
const now = Math.floor(Date.now() / 1000);

if (now - loginTime < 24 * 60 * 60) {
  return userInfo;
} else {
  clearUserInfo(); // 清除过期信息
}
```

## 🐛 常见问题

### 1. 跨域问题
确保后端正确设置CORS头：
```go
// 在 middleware/cors.go 中
c.Header("Access-Control-Allow-Origin", "https://carture.matrix-net.tech")
```

### 2. 重定向循环
确保URL清理逻辑正确执行：
```javascript
// 在授权成功后清理URL参数
const cleanURLParams = () => {
  const url = new URL(window.location);
  url.searchParams.delete('wechat_auth');
  url.searchParams.delete('user_info');
  window.history.replaceState({}, document.title, url.toString());
};
```

### 3. localStorage访问失败
在某些环境下localStorage可能不可用，增加异常处理：
```javascript
const saveUserInfo = (userInfo) => {
  try {
    if (typeof Storage !== 'undefined') {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(userInfo));
      return true;
    }
  } catch (error) {
    console.error('localStorage不可用:', error);
    // 可以使用其他存储方案，如React Native的AsyncStorage
  }
  return false;
};
```

## 📚 API参考

### useWechatAuth Hook

#### 参数
```javascript
const options = {
  autoCheck: true,           // 是否自动检查授权状态
  onAuthSuccess: (userInfo) => {}, // 授权成功回调
  onAuthError: (error) => {}       // 授权失败回调
};
```

#### 返回值
```javascript
const {
  userInfo,        // 用户信息对象
  isLoading,       // 是否正在加载
  error,           // 错误信息
  isAuthenticated, // 是否已认证
  isWechatBrowser, // 是否在微信浏览器中
  login,           // 登录函数
  logout,          // 登出函数
  refresh,         // 刷新用户信息
  checkAuth,       // 检查授权状态
  fetchUser,       // 获取指定用户信息
  clearError       // 清除错误
} = useWechatAuth(options);
```

### 用户信息结构
```javascript
{
  openid: "o3yLZ14BG3RwIQivHu9qWhs4b6gg",
  nickname: "tom",
  sex: 0,           // 0:未知, 1:男, 2:女
  language: "",
  city: "",
  province: "",
  country: "",
  headimgurl: "https://...",
  privilege: [],
  login_time: 1693920000  // Unix时间戳
}
```

## 🎉 完成！

现在你的React Expo应用已经支持微信授权登录了！用户在微信中访问你的应用时会自动检测登录状态，未登录时显示授权按钮，授权后保存用户信息到本地存储。
