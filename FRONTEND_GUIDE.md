# Vista 微信授权服务 - 前端集成指南

## 📋 概述

Vista 是一个基于 GraphQL 的微信公众号授权服务，提供微信 OAuth2 授权登录功能。前端可以通过 GraphQL API 实现微信用户授权和用户信息获取。

## 🚀 快速开始

### 服务地址
- **GraphQL Playground**: `http://localhost:8080/` (开发调试工具，浏览器访问)
- **GraphQL API**: `http://localhost:8080/wechat/query` (实际API端点，前端代码调用)

### 地址说明
- **Playground** 是一个可视化的开发工具，用于测试和调试 GraphQL 查询
- **API** 是真正的接口地址，前端应用通过 POST 请求调用这个地址
- 开发时在 Playground 测试查询，生产时前端代码调用 API 地址
- 使用 `/wechat` 路径前缀，便于区分不同业务模块

### GraphQL Schema

```graphql
type WechatAuth {
  openid: String!
  unionid: String
  accessToken: String!
  expiresIn: Int!
  refreshToken: String
  scope: String
  createdAt: String!
  updatedAt: String!
}

type Query {
  # 根据 openid 查询微信授权信息
  wechatAuth(openid: String!): WechatAuth
}

type Mutation {
  # 微信授权回调处理，用 code 换取用户信息
  wechatAuthCallback(code: String!): WechatAuth!
}
```

## 🔄 完整授权流程

### 1. 前端发起微信授权

在用户点击微信登录时，跳转到微信授权页面：

```javascript
function startWechatAuth() {
  const appId = 'wx103412583eab4dfd'; // 从服务端配置获取
  const redirectUri = encodeURIComponent('https://carture.matrix-net.tech/api/wechat/callback');
  const scope = 'snsapi_userinfo'; // 或 snsapi_base
  const state = Math.random().toString(36).substring(7); // 随机state值
  
  const authUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?` +
    `appid=${appId}` +
    `&redirect_uri=${redirectUri}` +
    `&response_type=code` +
    `&scope=${scope}` +
    `&state=${state}` +
    `#wechat_redirect`;
  
  // 保存 state 用于后续验证
  localStorage.setItem('wechat_auth_state', state);
  
  // 跳转到微信授权页面
  window.location.href = authUrl;
}
```

### 2. 处理微信回调

微信授权完成后，会回调到你配置的 `redirect_uri`，携带 `code` 参数：

```javascript
// 解析 URL 参数
function getUrlParams() {
  const params = new URLSearchParams(window.location.search);
  return {
    code: params.get('code'),
    state: params.get('state')
  };
}

// 处理微信回调
async function handleWechatCallback() {
  const { code, state } = getUrlParams();
  
  // 验证 state
  const savedState = localStorage.getItem('wechat_auth_state');
  if (state !== savedState) {
    throw new Error('Invalid state parameter');
  }
  
  if (!code) {
    throw new Error('授权失败：未获取到 code');
  }
  
  try {
    // 调用后端 GraphQL API 处理授权
    const result = await wechatAuthCallback(code);
    
    // 保存用户信息
    localStorage.setItem('wechat_user', JSON.stringify(result));
    
    // 跳转到主页或用户中心
    window.location.href = '/dashboard';
    
  } catch (error) {
    console.error('微信授权失败:', error);
    // 处理错误，比如跳转到错误页面
  }
}
```

### 3. GraphQL API 调用

#### 使用 Fetch API

```javascript
// 微信授权回调处理
async function wechatAuthCallback(code) {
  const query = `
    mutation WechatAuthCallback($code: String!) {
      wechatAuthCallback(code: $code) {
        openid
        unionid
        accessToken
        expiresIn
        refreshToken
        scope
        createdAt
        updatedAt
      }
    }
  `;
  
  const response = await fetch('http://localhost:8080/wechat/query', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query,
      variables: { code }
    })
  });
  
  const result = await response.json();
  
  if (result.errors) {
    throw new Error(result.errors[0].message);
  }
  
  return result.data.wechatAuthCallback;
}

// 查询用户授权信息
async function getWechatAuth(openid) {
  const query = `
    query GetWechatAuth($openid: String!) {
      wechatAuth(openid: $openid) {
        openid
        unionid
        accessToken
        expiresIn
        scope
        createdAt
        updatedAt
      }
    }
  `;
  
  const response = await fetch('http://localhost:8080/wechat/query', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query,
      variables: { openid }
    })
  });
  
  const result = await response.json();
  
  if (result.errors) {
    throw new Error(result.errors[0].message);
  }
  
  return result.data.wechatAuth;
}
```

#### 使用 Apollo Client

```javascript
import { ApolloClient, InMemoryCache, gql, useMutation, useQuery } from '@apollo/client';

// 创建 Apollo Client
const client = new ApolloClient({
  uri: 'http://localhost:8080/wechat/query',
  cache: new InMemoryCache()
});

// 定义 GraphQL 操作
const WECHAT_AUTH_CALLBACK = gql`
  mutation WechatAuthCallback($code: String!) {
    wechatAuthCallback(code: $code) {
      openid
      unionid
      accessToken
      expiresIn
      refreshToken
      scope
      createdAt
      updatedAt
    }
  }
`;

const GET_WECHAT_AUTH = gql`
  query GetWechatAuth($openid: String!) {
    wechatAuth(openid: $openid) {
      openid
      unionid
      accessToken
      expiresIn
      scope
      createdAt
      updatedAt
    }
  }
`;

// React 组件示例
function WechatLogin() {
  const [wechatAuthCallback] = useMutation(WECHAT_AUTH_CALLBACK);
  
  const handleCallback = async (code) => {
    try {
      const { data } = await wechatAuthCallback({
        variables: { code }
      });
      
      console.log('授权成功:', data.wechatAuthCallback);
      // 处理登录成功逻辑
      
    } catch (error) {
      console.error('授权失败:', error);
    }
  };
  
  return (
    <button onClick={startWechatAuth}>
      微信登录
    </button>
  );
}
```

## 🛠️ 实用工具函数

```javascript
// 微信授权工具类
class WechatAuth {
  constructor(config) {
    this.appId = config.appId;
    this.redirectUri = config.redirectUri;
    this.apiEndpoint = config.apiEndpoint || 'http://localhost:8080/wechat/query';
  }
  
  // 生成授权 URL
  generateAuthUrl(scope = 'snsapi_userinfo') {
    const state = this.generateState();
    this.saveState(state);
    
    return `https://open.weixin.qq.com/connect/oauth2/authorize?` +
      `appid=${this.appId}` +
      `&redirect_uri=${encodeURIComponent(this.redirectUri)}` +
      `&response_type=code` +
      `&scope=${scope}` +
      `&state=${state}` +
      `#wechat_redirect`;
  }
  
  // 发起授权
  authorize(scope = 'snsapi_userinfo') {
    window.location.href = this.generateAuthUrl(scope);
  }
  
  // 处理回调
  async handleCallback() {
    const { code, state } = this.getUrlParams();
    
    if (!this.validateState(state)) {
      throw new Error('Invalid state parameter');
    }
    
    if (!code) {
      throw new Error('Authorization failed: no code received');
    }
    
    return await this.exchangeCodeForToken(code);
  }
  
  // 用 code 换取 token
  async exchangeCodeForToken(code) {
    const query = `
      mutation WechatAuthCallback($code: String!) {
        wechatAuthCallback(code: $code) {
          openid
          unionid
          accessToken
          expiresIn
          refreshToken
          scope
          createdAt
          updatedAt
        }
      }
    `;
    
    const response = await fetch(this.apiEndpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query, variables: { code } })
    });
    
    const result = await response.json();
    
    if (result.errors) {
      throw new Error(result.errors[0].message);
    }
    
    return result.data.wechatAuthCallback;
  }
  
  // 工具方法
  generateState() {
    return Math.random().toString(36).substring(2, 15) +
           Math.random().toString(36).substring(2, 15);
  }
  
  saveState(state) {
    localStorage.setItem('wechat_auth_state', state);
  }
  
  validateState(state) {
    const savedState = localStorage.getItem('wechat_auth_state');
    localStorage.removeItem('wechat_auth_state');
    return state === savedState;
  }
  
  getUrlParams() {
    const params = new URLSearchParams(window.location.search);
    return {
      code: params.get('code'),
      state: params.get('state')
    };
  }
}

// 使用示例
const wechatAuth = new WechatAuth({
  appId: 'wx103412583eab4dfd',
  redirectUri: 'https://carture.matrix-net.tech/api/wechat/callback',
  apiEndpoint: 'http://localhost:8080/wechat/query'
});

// 发起授权
document.getElementById('wechat-login').onclick = () => {
  wechatAuth.authorize('snsapi_userinfo');
};

// 处理回调（在回调页面）
if (window.location.search.includes('code=')) {
  wechatAuth.handleCallback()
    .then(user => {
      console.log('登录成功:', user);
      // 保存用户信息并跳转
      localStorage.setItem('user', JSON.stringify(user));
      window.location.href = '/dashboard';
    })
    .catch(error => {
      console.error('登录失败:', error);
    });
}
```

## 📱 响应数据格式

### 成功响应
```json
{
  "data": {
    "wechatAuthCallback": {
      "openid": "oLVPpjqs9BhvzwPj5A-vTYAX3GLc",
      "unionid": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
      "accessToken": "ACCESS_TOKEN",
      "expiresIn": 7200,
      "refreshToken": "REFRESH_TOKEN",
      "scope": "snsapi_userinfo",
      "createdAt": "2025-09-04T19:30:00Z",
      "updatedAt": "2025-09-04T19:30:00Z"
    }
  }
}
```

### 错误响应
```json
{
  "errors": [
    {
      "message": "Invalid authorization code",
      "path": ["wechatAuthCallback"]
    }
  ]
}
```

## ⚠️ 注意事项

1. **安全性**
   - 使用 HTTPS 传输
   - 验证 `state` 参数防止 CSRF 攻击
   - 不要在前端存储敏感信息

2. **错误处理**
   - 处理网络错误
   - 处理授权被拒绝的情况
   - 提供用户友好的错误提示

3. **兼容性**
   - 确保在微信浏览器中正常工作
   - 测试不同版本的微信客户端

4. **用户体验**
   - 提供加载状态提示
   - 授权失败时给出明确提示
   - 支持重新授权

## 🔧 调试

### 使用 GraphQL Playground

访问 `http://localhost:8080/` 可以直接测试 GraphQL API：

```graphql
# 测试授权回调
mutation {
  wechatAuthCallback(code: "test_code_123") {
    openid
    accessToken
    expiresIn
  }
}

# 测试查询用户信息
query {
  wechatAuth(openid: "test_openid") {
    openid
    unionid
    scope
  }
}
```

### 开发环境配置

确保后端服务正在运行：
```bash
cd vista
go run main.go
```

服务启动后会显示：
```
Server starting on port 8080
GraphQL playground: http://localhost:8080/
GraphQL endpoint: http://localhost:8080/wechat/query
```

## 📞 技术支持

如有问题，请检查：
1. 后端服务是否正常运行
2. 微信公众号配置是否正确
3. 网络连接是否正常
4. GraphQL 查询语法是否正确
