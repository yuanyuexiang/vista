# Vista - 微信授权服务

简洁的微信OAuth授权服务，采用前端驱动模式。

## 🚀 特性

- **前端驱动** - 完全由前端控制授权流程
- **API简洁** - 仅提供2个核心API接口
- **智能检测** - 自动判断授权状态和过期时间
- **易于集成** - 适配任何前端框架

## 📋 API接口

### 1. 通过授权码获取用户信息
```
POST /api/wechat/auth
Content-Type: application/json

{
  "code": "微信授权码",
  "state": "可选状态参数"
}
```

### 2. 查询用户授权状态
```
GET /api/user/{openid}
```

## 🔧 快速开始

### 1. 配置
复制并修改配置文件：
```bash
cp config.yaml.example config.yaml
```

配置微信测试号信息：
```yaml
wechat:
  app_id: "你的微信AppID"
  app_secret: "你的微信AppSecret"
  redirect_uri: "https://你的域名/wechat/callback"
```

### 2. 启动服务
```bash
go run main.go
```

### 3. 访问演示
在微信中打开：`https://你的域名/wechat/frontend-driven`

## 📖 前端集成

### JavaScript示例
```javascript
// 1. 构建授权链接
function buildWechatAuthURL() {
    const appId = 'wx1eb05232cfbb49f7';
    const redirectURI = encodeURIComponent(window.location.href);
    const state = 'state_' + Date.now();
    
    return `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appId}&redirect_uri=${redirectURI}&response_type=code&scope=snsapi_userinfo&state=${state}#wechat_redirect`;
}

// 2. 处理微信回调
async function handleWechatCallback(code, state) {
    const response = await fetch('/api/wechat/auth', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code, state })
    });
    
    const data = await response.json();
    if (data.code === 200) {
        // 授权成功，保存用户信息
        localStorage.setItem('user_info', JSON.stringify(data.data.user_info));
    }
}

// 3. 检查授权状态
async function checkAuthStatus(openid) {
    const response = await fetch(`/api/user/${openid}`);
    const data = await response.json();
    
    if (data.data.need_auth) {
        // 需要重新授权
        window.location.href = buildWechatAuthURL();
    }
}
```

## 🏗️ 项目结构

```
vista/
├── main.go                           # 应用入口
├── config.yaml                       # 配置文件
├── wechat_frontend_driven.html       # 演示页面
├── internal/
│   ├── handler/
│   │   └── wechat.go                 # API处理器
│   ├── router/
│   │   └── router.go                 # 路由配置
│   ├── service/
│   │   └── wechat-oauth.go           # 微信OAuth服务
│   ├── model/
│   │   └── wechat_user.go            # 用户模型
│   └── repository/
│       └── wechat_user_repository.go # 数据访问层
└── pkg/
    ├── response/
    │   └── response.go               # 统一响应格式
    └── utils/
        └── wechat.go                 # 工具函数
```

## 🔄 授权流程

1. **前端构建授权链接** - 直接构建微信OAuth URL
2. **用户授权** - 微信返回授权码到回调页面
3. **前端获取code** - 从URL参数中提取授权码
4. **调用后端API** - POST /api/wechat/auth 获取用户信息
5. **本地存储** - 保存用户信息到localStorage
6. **状态检查** - 定期检查授权是否过期

## ⚙️ 配置说明

### 微信测试号
- 申请地址：https://developers.weixin.qq.com/sandbox
- 配置JS接口安全域名
- 配置授权回调域名

### 数据库
支持MySQL，自动创建表结构。

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License
