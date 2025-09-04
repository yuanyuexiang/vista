# Vista 数据库设计文档

## 📋 概述

Vista 微信授权服务的数据库设计，使用 MySQL 作为主数据库，通过 GORM 进行操作。

## 🗄️ 数据库表结构

### wechat_users 表

微信用户信息表，存储授权用户的基本信息和令牌。

```sql
CREATE TABLE `wechat_users` (
  `id` bigint unsigned AUTO_INCREMENT,
  `open_id` varchar(128) NOT NULL COMMENT '微信openid',
  `union_id` varchar(128) COMMENT '微信unionid',
  `nickname` varchar(100) COMMENT '用户昵称',
  `avatar` text COMMENT '头像URL',
  `sex` tinyint COMMENT '性别 0-未知 1-男 2-女',
  `province` varchar(50) COMMENT '省份',
  `city` varchar(50) COMMENT '城市',
  `country` varchar(50) COMMENT '国家',
  `access_token` text COMMENT '访问令牌',
  `refresh_token` text COMMENT '刷新令牌',
  `expires_at` datetime(3) NULL COMMENT '令牌过期时间',
  `scope` varchar(100) COMMENT '授权范围',
  `is_active` boolean DEFAULT true COMMENT '是否有效',
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_wechat_users_open_id` (`open_id`),
  INDEX `idx_wechat_users_union_id` (`union_id`)
);
```

### 字段说明

| 字段名 | 类型 | 说明 | 必填 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID，自增 | ✅ |
| open_id | varchar(128) | 微信用户唯一标识，每个应用下唯一 | ✅ |
| union_id | varchar(128) | 微信开放平台统一标识，跨应用唯一 | ❌ |
| nickname | varchar(100) | 用户昵称 | ❌ |
| avatar | text | 用户头像URL | ❌ |
| sex | tinyint | 性别：0-未知，1-男，2-女 | ❌ |
| province | varchar(50) | 省份 | ❌ |
| city | varchar(50) | 城市 | ❌ |
| country | varchar(50) | 国家 | ❌ |
| access_token | text | 微信访问令牌 | ❌ |
| refresh_token | text | 微信刷新令牌 | ❌ |
| expires_at | datetime(3) | 令牌过期时间 | ❌ |
| scope | varchar(100) | 授权范围：snsapi_base 或 snsapi_userinfo | ❌ |
| is_active | boolean | 用户是否有效，默认 true | ❌ |
| created_at | datetime(3) | 创建时间 | ❌ |
| updated_at | datetime(3) | 更新时间 | ❌ |

### 索引设计

- **主键索引**: `PRIMARY KEY (id)`
- **唯一索引**: `idx_wechat_users_open_id (open_id)` - 保证每个应用下 openid 唯一
- **普通索引**: `idx_wechat_users_union_id (union_id)` - 便于通过 unionid 查询

## 🔧 数据库功能

### 连接管理

```go
// 初始化数据库连接
func Init() error

// 获取数据库实例
func GetDB() *gorm.DB

// 数据库健康检查
func HealthCheck() error

// 关闭数据库连接
func Close() error
```

### 连接池配置

- **最大空闲连接数**: 10
- **最大打开连接数**: 100
- **连接最大生存时间**: 1小时

### 用户操作

```go
// 保存或更新微信用户信息
func SaveOrUpdateWechatUser(user *WechatUser) error

// 根据 openid 获取微信用户信息
func GetWechatUserByOpenID(openid string) (*WechatUser, error)
```

## 🚀 自动迁移

系统启动时会自动执行数据库迁移，创建或更新表结构：

```go
func AutoMigrate() error {
    return db.AutoMigrate(&WechatUser{})
}
```

## 🔍 使用示例

### 保存用户信息

```go
user := &database.WechatUser{
    OpenID:       "wx_openid_123",
    UnionID:      "wx_unionid_456",
    Nickname:     "张三",
    Avatar:       "https://wx.qlogo.cn/...",
    Sex:          1,
    Province:     "广东",
    City:         "深圳",
    Country:      "中国",
    AccessToken:  "access_token_xyz",
    RefreshToken: "refresh_token_abc",
    ExpiresAt:    time.Now().Add(2 * time.Hour),
    Scope:        "snsapi_userinfo",
    IsActive:     true,
}

err := database.SaveOrUpdateWechatUser(user)
```

### 查询用户信息

```go
user, err := database.GetWechatUserByOpenID("wx_openid_123")
if err != nil {
    log.Printf("User not found: %v", err)
} else {
    log.Printf("User: %s", user.Nickname)
}
```

### 健康检查

```go
if err := database.HealthCheck(); err != nil {
    log.Printf("Database health check failed: %v", err)
} else {
    log.Println("Database is healthy")
}
```

## ⚠️ 注意事项

### 数据安全

- `access_token` 和 `refresh_token` 使用 TEXT 类型存储
- 生产环境建议对敏感字段进行加密
- 定期清理过期的令牌

### 性能优化

- `open_id` 字段设置了唯一索引，查询性能优秀
- `union_id` 字段设置了普通索引，支持跨应用查询
- 使用连接池避免频繁创建连接

### 扩展性

- 表结构支持微信授权的所有标准字段
- `is_active` 字段支持软删除功能
- 时间字段使用毫秒精度，便于精确控制

### 兼容性

- 支持 GORM v1.30.2+
- 兼容 MySQL 5.7+
- 字符集使用 utf8mb4，支持 emoji 等特殊字符

## 🔄 数据流程

1. **用户授权**: 用户通过微信授权，系统获取 code
2. **换取令牌**: 后端用 code 换取 access_token 和用户信息
3. **保存用户**: 调用 `SaveOrUpdateWechatUser` 保存到数据库
4. **用户查询**: 通过 `GetWechatUserByOpenID` 查询用户信息
5. **令牌刷新**: 使用 refresh_token 更新 access_token

## 📊 监控建议

- 监控数据库连接数和响应时间
- 定期检查过期令牌数量
- 监控用户增长趋势
- 记录数据库操作日志
