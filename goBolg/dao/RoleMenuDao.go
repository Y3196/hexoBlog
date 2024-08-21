package dao

import (
	"context"
	"goBolg/model"
	"gorm.io/gorm"
)

// RoleMenuDao 角色菜单 DAO 接口
type RoleMenuDao interface {
	// 查询是否有角色关联
	CountByMenuID(ctx context.Context, menuID uint) (int64, error)

	RemoveByRoleID(ctx context.Context, roleID int) error

	SaveBatch(ctx context.Context, roleMenus []model.RoleMenu) error
}

type roleMenuDao struct {
	db *gorm.DB
}

// NewRoleMenuDao 创建新的 RoleMenuDao 实例
func NewRoleMenuDao(db *gorm.DB) RoleMenuDao {
	return &roleMenuDao{db: db}
}

// CountByMenuID 查询是否有角色关联
func (dao *roleMenuDao) CountByMenuID(ctx context.Context, menuID uint) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).
		Model(&model.RoleMenu{}).
		Where("menu_id = ?", menuID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 在 dao/role_menu_dao.go 中
func (dao *roleMenuDao) RemoveByRoleID(ctx context.Context, roleID int) error {
	return dao.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&model.RoleMenu{}).Error
}

// 在 dao/role_menu_dao.go 中
func (dao *roleMenuDao) SaveBatch(ctx context.Context, roleMenus []model.RoleMenu) error {
	return dao.db.WithContext(ctx).
		Create(&roleMenus).Error
}
