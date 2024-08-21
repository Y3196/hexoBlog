package dao

import (
	"context"
	"goBolg/model"
	"gorm.io/gorm"
)

// RoleResourceDao 角色资源 DAO 接口
type RoleResourceDao interface {
	// 查询是否有角色关联
	CountByResourceID(ctx context.Context, resourceID uint) (int64, error)

	RemoveByRoleID(ctx context.Context, roleID int) error

	SaveBatch(ctx context.Context, roleResources []model.RoleResource) error
}

type roleResourceDao struct {
	db *gorm.DB
}

// NewRoleResourceDao 创建新的 RoleResourceDao 实例
func NewRoleResourceDao(db *gorm.DB) RoleResourceDao {
	return &roleResourceDao{db: db}
}

// CountByResourceID 查询是否有角色关联
func (dao *roleResourceDao) CountByResourceID(ctx context.Context, resourceID uint) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).
		Model(&model.RoleResource{}).
		Where("resource_id = ?", resourceID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 在 dao/role_resource_dao.go 中
func (dao *roleResourceDao) RemoveByRoleID(ctx context.Context, roleID int) error {
	return dao.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&model.RoleResource{}).Error
}

// 在 dao/role_resource_dao.go 中
func (dao *roleResourceDao) SaveBatch(ctx context.Context, roleResources []model.RoleResource) error {
	return dao.db.WithContext(ctx).
		Create(&roleResources).Error
}
