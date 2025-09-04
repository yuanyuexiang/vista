# Vista 微信授权服务 - 前端集成指南

## 📋 概述

基于 GraphQL 的微信公众号授权服务，提供微信 OAuth2 授权登录功能。

## 🚀 服务地址

- **GraphQL Playground**: `http://localhost:8080/` (开发调试)
- **GraphQL API**: `http://localhost:8080/wechat/query` (接口调用)

## 📡 GraphQL Schema

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
  wechatAuth(openid: String!): WechatAuth
}

type Mutation {
  wechatAuthCallback(code: String!): WechatAuth!
}
```

## 🔄 授权流程

### 1. 发起微信授权

```javascript
function startWechatAuth() {
  const appId = 'wx103412583eab4dfd';
  const redirectUri = 'https://carture.matrix-net.tech/api/wechat/callback';
  const state = Math.random().toString(36).substring(7);
  
  localStorage.setItem('wechat_auth_state', state);
  
  const authUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?` +
    `appid=${appId}&redirect_uri=${encodeURIComponent(redirectUri)}` +
    `&response_type=code&scope=snsapi_userinfo&state=${state}#wechat_redirect`;
  
  window.location.href = authUrl;
}
```

### 2. 处理回调

```javascript
async function handleWechatCallback() {
  const params = new URLSearchParams(window.location.search);
  const code = params.get('code');
  const state = params.get('state');
  
  // 验证 state
  if (state !== localStorage.getItem('wechat_auth_state')) {
    throw new Error('Invalid state');
  }
  
  // 调用 GraphQL API
  const result = await wechatAuthCallback(code);
  localStorage.setItem('wechat_user', JSON.stringify(result));
  
  window.location.href = '/dashboard';
}
```

### 3. GraphQL 调用

```javascript
async function wechatAuthCallback(code) {
  const query = `
    mutation WechatAuthCallback($code: String!) {
      wechatAuthCallback(code: $code) {
        openid
        accessToken
        expiresIn
      }
    }
  `;
  
  const response = await fetch('http://localhost:8080/wechat/query', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables: { code } })
  });
  
  const result = await response.json();
  return result.data.wechatAuthCallback;
}

async function getWechatAuth(openid) {
  const query = `
    query GetWechatAuth($openid: String!) {
      wechatAuth(openid: $openid) {
        openid
        accessToken
        scope
      }
    }
  `;
  
  const response = await fetch('http://localhost:8080/wechat/query', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables: { openid } })
  });
  
  const result = await response.json();
  return result.data.wechatAuth;
}
```

## 🔧 测试

在 GraphQL Playground (`http://localhost:8080/`) 中测试：

```graphql
# 测试授权回调
mutation {
  wechatAuthCallback(code: "test_code_123") {
    openid
    accessToken
  }
}

# 测试查询用户信息
query {
  wechatAuth(openid: "test_openid") {
    openid
    scope
  }
}
```

## ⚠️ 注意事项

- 使用 HTTPS 传输
- 验证 `state` 参数防止 CSRF
- 处理授权失败情况
- 确保微信浏览器兼容性
