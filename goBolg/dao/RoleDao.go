package dao

import (
	"context"
	"database/sql"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
)

// RoleDao 角色 DAO 接口
type RoleDao interface {
	// 查询路由角色列表
	ListResourceRoles(ctx context.Context) ([]dto.ResourceRoleDTO, error)

	// 根据用户 ID 获取角色列表
	ListRolesByUserInfoId(ctx context.Context, userInfoId int) ([]string, error)

	// 查询角色列表
	ListRoles(ctx context.Context, current, size int, conditionVO *vo.ConditionVO) ([]dto.RoleDTO, error)

	// 查询角色列表（带选择字段）
	ListRoleList(ctx context.Context) ([]model.Role, error)

	// 查询总量
	CountRoles(ctx context.Context, conditionVO *vo.ConditionVO) (int64, error)

	// 判断角色名重复
	CheckRoleExists(ctx context.Context, roleName string) (*model.Role, error)

	SaveOrUpdate(ctx context.Context, role *model.Role, tx *gorm.DB) error

	// 删除角色
	DeleteRoles(ctx context.Context, roleIdList []int) error
}

type roleDao struct {
	db *gorm.DB
}

// NewRoleDao 创建新的 RoleDao 实例
func NewRoleDao(db *gorm.DB) RoleDao {
	return &roleDao{db: db}
}

// ListResourceRoles 查询路由角色列表
func (dao *roleDao) ListResourceRoles(ctx context.Context) ([]dto.ResourceRoleDTO, error) {
	var roles []dto.ResourceRoleDTO
	err := dao.db.WithContext(ctx).
		Raw(`
			SELECT
				re.id,
				url,
				request_method,
				role_label
			FROM
				tb_resource re
				LEFT JOIN tb_role_resource rep ON re.id = rep.resource_id
				LEFT JOIN tb_role r ON rep.role_id = r.id
			WHERE
				parent_id IS NOT NULL
				AND is_anonymous = 0
		`).Scan(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// ListRolesByUserInfoId 根据用户 ID 获取角色列表
func (dao *roleDao) ListRolesByUserInfoId(ctx context.Context, userInfoId int) ([]string, error) {
	var roles []string
	err := dao.db.WithContext(ctx).
		Raw(`
			SELECT
				role_label
			FROM
				tb_role r
				JOIN tb_user_role ur ON r.id = ur.role_id
			WHERE
				ur.user_id = ?
		`, userInfoId).Scan(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// ListRoles 查询角色列表
func (dao *roleDao) ListRoles(ctx context.Context, current, size int, conditionVO *vo.ConditionVO) ([]dto.RoleDTO, error) {
	var roles []dto.RoleDTO
	query := dao.db.WithContext(ctx).
		Model(&model.Role{}).
		Select("tb_role.id, role_name, role_label, create_time, is_disable, GROUP_CONCAT(DISTINCT rr.resource_id) as resource_id_list, GROUP_CONCAT(DISTINCT rm.menu_id) as menu_id_list").
		Joins("LEFT JOIN tb_role_resource rr ON tb_role.id = rr.role_id").
		Joins("LEFT JOIN tb_role_menu rm ON tb_role.id = rm.role_id").
		Group("tb_role.id")

	if conditionVO.Keywords != nil && *conditionVO.Keywords != "" {
		query = query.Where("role_name LIKE ?", "%"+*conditionVO.Keywords+"%")
	}

	rows, err := query.Offset(current).Limit(size).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role dto.RoleDTO
		var resourceIDList sql.NullString
		var menuIDList sql.NullString

		err := rows.Scan(&role.ID, &role.RoleName, &role.RoleLabel, &role.CreateTime, &role.IsDisable, &resourceIDList, &menuIDList)
		if err != nil {
			return nil, err
		}

		if resourceIDList.Valid {
			role.ResourceIDList = dto.ParseConcatList(resourceIDList.String)
		} else {
			role.ResourceIDList = nil
		}

		if menuIDList.Valid {
			role.MenuIDList = dto.ParseConcatList(menuIDList.String)
		} else {
			role.MenuIDList = nil
		}

		roles = append(roles, role)
	}

	return roles, nil
}

// ListRoleList 查询角色列表（带选择字段）
func (dao *roleDao) ListRoleList(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	err := dao.db.WithContext(ctx).
		Model(&model.Role{}).
		Select("id, role_name").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// CountRoles 查询总量
func (dao *roleDao) CountRoles(ctx context.Context, conditionVO *vo.ConditionVO) (int64, error) {
	var count int64
	query := dao.db.WithContext(ctx).
		Model(&model.Role{})

	if conditionVO.Keywords != nil && *conditionVO.Keywords != "" {
		query = query.Where("role_name LIKE ?", "%"+*conditionVO.Keywords+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CheckRoleExists 判断角色名重复
func (dao *roleDao) CheckRoleExists(ctx context.Context, roleName string) (*model.Role, error) {
	var role model.Role
	err := dao.db.WithContext(ctx).
		Where("role_name = ?", roleName).
		First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &role, nil
}

func (dao *roleDao) SaveOrUpdate(ctx context.Context, role *model.Role, tx *gorm.DB) error {
	if role.ID == 0 {
		return tx.WithContext(ctx).Create(role).Error
	} else {
		return tx.WithContext(ctx).Save(role).Error
	}
}

// DeleteRoles 删除角色
func (dao *roleDao) DeleteRoles(ctx context.Context, roleIdList []int) error {
	return dao.db.WithContext(ctx).Where("id IN ?", roleIdList).Delete(&model.Role{}).Error
}
