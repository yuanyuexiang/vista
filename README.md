# Vista - 微信OAuth认证服务

Vista是一个简洁高效的微信OAuth认证服务，采用前端驱动模式，为开发者提供最简单的微信登录集成方案。

## 🎯 核心特性

- **🚀 前端驱动** - 授权流程完全由前端控制，灵活性更高
- **⚡ 极简API** - 只有2个核心接口，学习成本低
- **🔒 智能授权** - 自动判断授权状态，24小时有效期
- **📱 完美适配** - 专为微信内置浏览器优化
- **🛠️ 易于集成** - 支持原生JS、React、Vue等任何前端框架

## 🏗️ 系统架构

```
前端应用 ──────► 微信OAuth ──────► Vista后端
    │              │                │
    │ 1.构建授权链接  │                │
    │ ────────────► │                │
    │              │ 2.用户授权获取code │
    │ ◄──────────── │                │
    │ 3.发送code     │                │
    │ ──────────────────────────────► │
    │              │ 4.获取用户信息      │
    │              │ ◄──────────────── │
    │ 5.返回用户数据  │                │
    │ ◄──────────────────────────────── │
```

## 📋 快速开始

### 1. 启动服务

```bash
# 编译运行
go build -o vista .
./vista

# 服务运行在 https://carture.matrix-net.tech
```

### 2. 前端集成

#### 基础HTML示例

```html
<script>
const CONFIG = {
    API_BASE: 'https://carture.matrix-net.tech',
    WECHAT_APP_ID: 'wx1eb05232cfbb49f7'
};

// 构建授权链接
function buildWechatAuthURL() {
    const redirectURI = encodeURIComponent(window.location.href.split('?')[0]);
    const state = Date.now().toString();
    return `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${CONFIG.WECHAT_APP_ID}&redirect_uri=${redirectURI}&response_type=code&scope=snsapi_userinfo&state=${state}#wechat_redirect`;
}

// 开始授权
function startAuth() {
    window.location.href = buildWechatAuthURL();
}

// 处理回调
async function handleCallback(code) {
    const response = await fetch(`${CONFIG.API_BASE}/vista/wechat/api/auth`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code })
    });
    
    const data = await response.json();
    if (data.code === 200) {
        console.log('用户信息:', data.data.user_info);
    }
}
</script>
```

#### React Hook示例

```javascript
import { useState, useEffect } from 'react';

const useWechatAuth = () => {
    const [user, setUser] = useState(null);
    const [needAuth, setNeedAuth] = useState(false);
    
    // 检查授权状态
    const checkAuth = async () => {
        const openid = localStorage.getItem('wechat_openid');
        if (!openid) {
            setNeedAuth(true);
            return;
        }
        
        const response = await fetch(`https://carture.matrix-net.tech/vista/wechat/api/user/${openid}`);
        const data = await response.json();
        
        if (data.code === 200 && !data.data.need_auth) {
            setUser(data.data.user_info);
        } else {
            setNeedAuth(true);
        }
    };
    
    return { user, needAuth, checkAuth };
};
```

## 🔗 API接口

### 1. 授权接口
```
POST /vista/wechat/api/auth
```

**请求体：**
```json
{
  "code": "微信授权码",
  "state": "状态参数(可选)"
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
      "headimgurl": "头像URL",
      "sex": 1,
      "province": "省份",
      "city": "城市",
      "country": "国家",
      "expires_at": 1725552000
    }
  }
}
```

### 2. 用户信息查询
```
GET /vista/wechat/api/user/{openid}
```

**响应：**
```json
{
  "code": 200,
  "data": {
    "exists": true,
    "need_auth": false,
    "user_info": { /* 用户详细信息 */ }
  }
}
```

## � 项目结构

```
vista/
├── main.go                 # 程序入口
├── config.yaml             # 配置文件  
├── internal/
│   ├── handler/            # HTTP处理器
│   ├── service/            # 业务逻辑
│   ├── repository/         # 数据访问
│   ├── model/              # 数据模型
│   └── router/             # 路由配置
├── pkg/
│   ├── logger/             # 日志工具
│   ├── response/           # 响应工具
│   └── utils/              # 工具函数
└── wechat_frontend_driven.html  # 演示页面
```

## 🔧 配置说明

### config.yaml
```yaml
server:
  port: 8080

database:
  driver: mysql
  dsn: "user:password@tcp(localhost:3306)/vista?charset=utf8mb4&parseTime=True&loc=Local"

wechat:
  app_id: "wx1eb05232cfbb49f7"
  app_secret: "your_app_secret"

log:
  level: "info"
  file_path: "./server.log"
```

## 🧪 测试体验

访问在线演示页面体验完整流程：
```
https://carture.matrix-net.tech/vista/wechat/frontend-driven
```

## 📚 文档

- **[前端集成指南](./FRONTEND_INTEGRATION_GUIDE.md)** - 详细的前端集成文档，包含完整示例
- **[API文档](./API_DOCUMENTATION.md)** - 完整的API接口说明
- **[架构文档](./ARCHITECTURE.md)** - 系统架构设计说明

## 🛠️ 开发工具

- **Go 1.21+** - 后端开发语言
- **Gin** - Web框架
- **GORM** - ORM工具
- **MySQL** - 数据库

## 🎨 特色功能

- ✅ **智能授权检查** - 自动判断是否需要重新授权
- ✅ **24小时有效期** - 合理的授权时间窗口
- ✅ **前端驱动模式** - 给前端更多控制权
- ✅ **简洁API设计** - 只有2个核心接口
- ✅ **完整错误处理** - 友好的错误信息和状态码
- ✅ **本地存储支持** - 合理利用localStorage
- ✅ **移动端优化** - 专为微信浏览器设计

## 🚀 快速集成

对于前端开发者，我们提供了：

1. **📖 [完整集成指南](./FRONTEND_INTEGRATION_GUIDE.md)** - 包含详细步骤和代码示例
2. **🎯 React Hook示例** - 开箱即用的React集成方案
3. **📱 原生JS示例** - 适用于任何前端项目
4. **🔧 Vue.js示例** - Vue开发者的最佳实践
5. **💡 最佳实践** - 性能优化和用户体验建议

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个项目！

## 📄 许可证

MIT License

---

*简洁、高效、易用 - Vista让微信OAuth集成变得简单*
