package dao

import (
	"context"
	"goBolg/model"
	"gorm.io/gorm"
)

// UserRoleDao 接口
type UserRoleDao interface {
	CountRolesWithUsers(ctx context.Context, roleIdList []int) (int64, error)

	InsertUserRole(ctx context.Context, userRole *model.UserRole) error
}

// userRoleDaoImpl 实现 UserRoleDao 接口
type userRoleDaoImpl struct {
	db *gorm.DB
}

// NewUserRoleDao 创建新的 UserRoleDao 实例
func NewUserRoleDao(db *gorm.DB) UserRoleDao {
	return &userRoleDaoImpl{db: db}
}

// CountRolesWithUsers 判断角色下是否有用户
func (dao *userRoleDaoImpl) CountRolesWithUsers(ctx context.Context, roleIdList []int) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.UserRole{}).
		Where("role_id IN ?", roleIdList).
		Count(&count).Error
	return count, err
}

// InsertUserRole 插入用户角色关系
func (dao *userRoleDaoImpl) InsertUserRole(ctx context.Context, userRole *model.UserRole) error {
	return dao.db.WithContext(ctx).Create(userRole).Error
}
