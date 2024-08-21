package dao

import (
	"context"
	"fmt"
	"goBolg/model"
	"gorm.io/gorm"
	"time"
)

// UserInfoDao 接口
type UserInfoDao interface {
	CountUsers(ctx context.Context) (int64, error)

	GetEmailById(ctx context.Context, userId uint) (string, error)

	InsertUserInfo(ctx context.Context, userInfo *model.UserInfo) error

	GetUserInfoById(ctx context.Context, userId int) (*model.UserInfo, error)

	UpdateUserInfo(ctx context.Context, userInfo *model.UserInfo) error

	// 新增：根据用户名查询用户信息
	GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error)
}

// 实现
type userInfoDaoImpl struct {
	db *gorm.DB
}

func NewUserInfoDao(db *gorm.DB) UserInfoDao {
	return &userInfoDaoImpl{db: db}
}

func (dao *userInfoDaoImpl) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	err := dao.db.Table("tb_user_info").Count(&count).Error
	return count, err
}

// GetEmailById 根据 ID 查询邮箱
func (dao *userInfoDaoImpl) GetEmailById(ctx context.Context, userId uint) (string, error) {
	var user model.UserInfo
	err := dao.db.WithContext(ctx).Select("email").First(&user, userId).Error
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

// InsertUserInfo 插入用户信息
func (dao *userInfoDaoImpl) InsertUserInfo(ctx context.Context, userInfo *model.UserInfo) error {
	return dao.db.WithContext(ctx).Create(userInfo).Error
}

func (dao *userInfoDaoImpl) GetUserInfoById(ctx context.Context, userId int) (*model.UserInfo, error) {
	// 检查 userId 是否有效
	if userId == 0 {
		return nil, fmt.Errorf("invalid userId: %d", userId)
	}

	var user model.UserInfo
	err := dao.db.WithContext(ctx).First(&user, userId).Error
	if err != nil {
		return nil, fmt.Errorf("record not found for userId: %d, error: %w", userId, err)
	}
	return &user, nil
}

// UpdateUserInfo 更新用户信息
func (dao *userInfoDaoImpl) UpdateUserInfo(ctx context.Context, userInfo *model.UserInfo) error {
	userInfo.UpdateTime = time.Now() // 设置当前时间
	// 仅更新除 create_time 之外的字段
	return dao.db.WithContext(ctx).Model(&userInfo).Omit("CreateTime").Updates(userInfo).Error
}

// 新增：根据用户名查询用户信息
func (dao *userInfoDaoImpl) GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error) {
	var user model.UserInfo
	err := dao.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		// 如果找不到记录，返回 nil 和错误
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
