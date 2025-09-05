package repository

import (
	"fmt"
	"vista/internal/model"

	"gorm.io/gorm"
)

// WechatUserRepository 微信用户数据访问层
type WechatUserRepository struct {
	db *gorm.DB
}

// NewWechatUserRepository 创建微信用户仓库实例
func NewWechatUserRepository(db *gorm.DB) *WechatUserRepository {
	return &WechatUserRepository{db: db}
}

// Save 保存或更新微信用户信息
func (r *WechatUserRepository) Save(user *model.WechatUser) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// 根据 openid 查找用户，如果存在则更新，不存在则创建
	var existingUser model.WechatUser
	result := r.db.Where("openid = ?", user.OpenID).First(&existingUser)

	if result.Error == gorm.ErrRecordNotFound {
		// 用户不存在，创建新用户
		return r.db.Create(user).Error
	} else if result.Error != nil {
		// 查询出错
		return result.Error
	}

	// 用户存在，更新信息
	user.ID = existingUser.ID
	user.CreatedAt = existingUser.CreatedAt
	return r.db.Save(user).Error
}

// GetByOpenID 根据 openid 获取微信用户信息
func (r *WechatUserRepository) GetByOpenID(openid string) (*model.WechatUser, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var user model.WechatUser
	err := r.db.Where("openid = ? AND is_active = ?", openid, true).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID 根据 ID 获取用户信息
func (r *WechatUserRepository) GetByID(id uint) (*model.WechatUser, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var user model.WechatUser
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAll 获取所有活跃用户
func (r *WechatUserRepository) GetAll() ([]model.WechatUser, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var users []model.WechatUser
	err := r.db.Where("is_active = ?", true).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Delete 软删除用户
func (r *WechatUserRepository) Delete(openid string) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	return r.db.Model(&model.WechatUser{}).
		Where("openid = ?", openid).
		Update("is_active", false).Error
}
