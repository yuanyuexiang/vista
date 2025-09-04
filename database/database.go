package database

import (
	"fmt"
	"vista/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

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

	// 自动迁移
	// if err := AutoMigrate(); err != nil {
	// 	return fmt.Errorf("failed to migrate database: %v", err)
	// }

	return nil
}

// // AutoMigrate 自动迁移数据库表
// func AutoMigrate() error {
// 	return db.AutoMigrate(
// 		&model.OAuthToken{},
// 		&model.Keyword{},
// 		&model.Email{},
// 		&model.Activity{},
// 	)
// }

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
