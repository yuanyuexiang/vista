# Vista 项目架构文档

## 📋 项目概述

Vista 是一个基于 Go 语言和 Gin 框架的微信公众号授权服务，采用现代化的项目架构和最佳实践。

## 🏗️ 项目结构

```
vista/
├── cmd/                    # 应用程序入口点（暂未使用，保留扩展）
├── config/                 # 配置管理
│   ├── config.go          # 配置结构和初始化
│   └── config.yaml        # 配置文件
├── internal/              # 私有应用代码
│   ├── database/          # 数据库连接和管理
│   │   └── database.go
│   ├── handler/           # HTTP 处理器层
│   │   └── wechat.go      # 微信相关接口
│   ├── middleware/        # 中间件
│   │   └── middleware.go  # 日志、恢复、CORS 中间件
│   ├── model/             # 数据模型层
│   │   └── wechat_user.go # 微信用户模型
│   ├── repository/        # 数据访问层
│   │   └── wechat_user_repository.go
│   ├── router/            # 路由配置
│   │   └── router.go
│   └── service/           # 业务逻辑层
│       └── wechat-oauth.go # 微信 OAuth 服务
├── pkg/                   # 可重用的公共包
│   ├── logger/            # 日志工具
│   │   └── logger.go
│   ├── response/          # 统一响应格式
│   │   └── response.go
│   └── utils/             # 工具函数
│       └── wechat.go      # 微信相关工具
├── docs/                  # 文档
│   ├── DATABASE.md        # 数据库设计文档
│   └── FRONTEND_GUIDE.md  # 前端集成指南
├── main.go               # 程序入口
├── go.mod               # 依赖管理
├── go.sum
├── Dockerfile           # Docker 构建文件
├── docker-compose.yaml  # Docker Compose 配置
└── README.md           # 项目说明
```

## 🎯 架构设计原则

### 1. 分层架构

采用经典的分层架构模式：

```
┌─────────────────┐
│   Handler 层     │ ← HTTP 请求处理
├─────────────────┤
│   Service 层     │ ← 业务逻辑
├─────────────────┤
│  Repository 层   │ ← 数据访问
├─────────────────┤
│   Model 层      │ ← 数据模型
└─────────────────┘
```

### 2. 依赖注入

使用构造函数注入，实现松耦合：

```go
// Service 依赖 Repository
func NewWechatOAuthService(userRepo *repository.WechatUserRepository) *WechatOAuthService

// Repository 依赖 Database
func NewWechatUserRepository(db *gorm.DB) *WechatUserRepository
```

### 3. 接口隔离

每层都定义了清晰的接口，便于测试和扩展。

## 📦 各层职责

### Handler 层 (`internal/handler/`)

- **职责**: HTTP 请求处理、参数验证、响应格式化
- **依赖**: Service 层
- **特点**: 
  - 不包含业务逻辑
  - 使用统一的响应格式
  - 处理 HTTP 状态码

### Service 层 (`internal/service/`)

- **职责**: 业务逻辑实现、第三方 API 调用
- **依赖**: Repository 层
- **特点**:
  - 无状态设计
  - 包含核心业务规则
  - 处理外部 API 交互

### Repository 层 (`internal/repository/`)

- **职责**: 数据访问、CRUD 操作
- **依赖**: Model 层、Database
- **特点**:
  - 封装数据库操作
  - 提供领域对象的持久化
  - 隔离数据访问逻辑

### Model 层 (`internal/model/`)

- **职责**: 数据结构定义、领域对象
- **依赖**: 无
- **特点**:
  - 纯数据结构
  - 包含领域规则
  - GORM 标签配置

## 🔧 技术栈

### 核心框架
- **Web 框架**: Gin v1.10.1
- **ORM**: GORM v1.30.2
- **数据库**: MySQL 5.7+
- **配置管理**: Viper v1.20.1

### 开发工具
- **Go 版本**: 1.24+
- **容器化**: Docker + Docker Compose
- **依赖管理**: Go Modules

## 🚀 最佳实践

### 1. 代码组织

- ✅ 使用 `internal/` 目录保护私有代码
- ✅ 使用 `pkg/` 目录存放可复用包
- ✅ 按功能模块组织代码
- ✅ 清晰的包命名和文件命名

### 2. 错误处理

```go
// 统一错误响应
if err != nil {
    response.InternalServerError(c, fmt.Sprintf("operation failed: %v", err))
    return
}
```

### 3. 日志记录

```go
// 结构化日志
logger.Infof("User saved successfully - OpenID: %s", openID)
logger.Errorf("Database operation failed: %v", err)
```

### 4. 依赖注入

```go
// 在 Handler 中创建依赖
userRepo := repository.NewWechatUserRepository(database.GetDB())
wechatService := service.NewWechatOAuthService(userRepo)
```

### 5. 中间件使用

```go
// 路由级中间件
r.Use(middleware.Logger())
r.Use(middleware.Recovery())
r.Use(middleware.CORS())
```

## 🔍 设计模式

### 1. Repository 模式
- 数据访问层抽象
- 便于单元测试
- 数据库无关性

### 2. Service 模式
- 业务逻辑封装
- 事务管理
- 跨领域操作

### 3. 依赖注入模式
- 松耦合设计
- 便于测试
- 配置灵活

## 📊 性能优化

### 1. 数据库连接池

```go
sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间
```

### 2. 索引优化

```sql
UNIQUE INDEX `idx_wechat_users_open_id` (`open_id`)
INDEX `idx_wechat_users_union_id` (`union_id`)
```

### 3. 响应缓存

预留缓存层接口，便于后续添加 Redis 缓存。

## 🧪 测试策略

### 1. 单元测试
- Service 层业务逻辑测试
- Repository 层数据访问测试
- 工具函数测试

### 2. 集成测试
- API 端到端测试
- 数据库集成测试

### 3. Mock 测试
- 外部 API 调用 Mock
- 数据库操作 Mock

## 🔒 安全考虑

### 1. 输入验证
- 参数格式验证
- SQL 注入防护
- XSS 防护

### 2. 状态管理
- CSRF 防护（state 参数）
- 会话管理
- 令牌安全存储

### 3. 错误处理
- 敏感信息隐藏
- 统一错误响应
- 详细错误日志

## 📈 扩展性设计

### 1. 水平扩展
- 无状态设计
- 数据库读写分离就绪
- 负载均衡友好

### 2. 功能扩展
- 插件化中间件
- 多种授权方式支持
- 微服务拆分就绪

### 3. 配置扩展
- 多环境配置
- 动态配置更新
- 功能开关

## 🎯 总结

当前项目结构完全符合 Go 语言的最佳实践：

- ✅ **标准目录布局**: 遵循 Go 社区标准
- ✅ **清晰的分层**: 职责分离，易于维护
- ✅ **依赖注入**: 松耦合，便于测试
- ✅ **错误处理**: 统一的错误处理机制
- ✅ **可扩展性**: 支持水平和垂直扩展
- ✅ **安全性**: 多层安全防护
- ✅ **性能**: 连接池和索引优化
- ✅ **可维护性**: 清晰的代码组织

这个架构为项目的长期发展奠定了坚实的基础。
