package database

import (
	"fmt"
	"log"
	"time"
	"vista/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// WechatUser 微信用户模型
type WechatUser struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	OpenID       string    `json:"openid" gorm:"type:varchar(128);uniqueIndex;not null;comment:微信openid"`
	UnionID      string    `json:"unionid" gorm:"type:varchar(128);index;comment:微信unionid"`
	Nickname     string    `json:"nickname" gorm:"type:varchar(100);comment:用户昵称"`
	Avatar       string    `json:"avatar" gorm:"type:text;comment:头像URL"`
	Sex          int       `json:"sex" gorm:"type:tinyint;comment:性别 0-未知 1-男 2-女"`
	Province     string    `json:"province" gorm:"type:varchar(50);comment:省份"`
	City         string    `json:"city" gorm:"type:varchar(50);comment:城市"`
	Country      string    `json:"country" gorm:"type:varchar(50);comment:国家"`
	AccessToken  string    `json:"access_token" gorm:"type:text;comment:访问令牌"`
	RefreshToken string    `json:"refresh_token" gorm:"type:text;comment:刷新令牌"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"comment:令牌过期时间"`
	Scope        string    `json:"scope" gorm:"type:varchar(100);comment:授权范围"`
	IsActive     bool      `json:"is_active" gorm:"default:true;comment:是否有效"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Init 初始化数据库连接
func Init() error {
	cfg := config.Get()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// 自动迁移
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	return db.AutoMigrate(
		&WechatUser{},
	)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

// HealthCheck 数据库健康检查
func HealthCheck() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %v", err)
	}

	return sqlDB.Ping()
}

// Close 关闭数据库连接
func Close() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// SaveOrUpdateWechatUser 保存或更新微信用户信息
func SaveOrUpdateWechatUser(user *WechatUser) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// 根据 openid 查找用户，如果存在则更新，不存在则创建
	var existingUser WechatUser
	result := db.Where("openid = ?", user.OpenID).First(&existingUser)

	if result.Error == gorm.ErrRecordNotFound {
		// 用户不存在，创建新用户
		return db.Create(user).Error
	} else if result.Error != nil {
		// 查询出错
		return result.Error
	}

	// 用户存在，更新信息
	user.ID = existingUser.ID
	user.CreatedAt = existingUser.CreatedAt
	return db.Save(user).Error
}

// GetWechatUserByOpenID 根据 openid 获取微信用户信息
func GetWechatUserByOpenID(openid string) (*WechatUser, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var user WechatUser
	err := db.Where("openid = ? AND is_active = ?", openid, true).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
