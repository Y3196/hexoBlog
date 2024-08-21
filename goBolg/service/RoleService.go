package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

// RoleService 角色服务接口
type RoleService interface {
	// 查询用户角色列表
	ListUserRoles(ctx context.Context) ([]dto.UserRoleDTO, error)

	ListRoles(ctx context.Context, conditionVO *vo.ConditionVO) (*vo.PageResult[dto.RoleDTO], error)

	DeleteRoles(ctx context.Context, roleIdList []int) error

	SaveOrUpdateRole(ctx context.Context, roleVO *vo.RoleVO) error
}
