package Impl

import (
	"context"
	"errors"
	"fmt"
	constants "goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/handler"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"gorm.io/gorm"
)

// roleServiceImpl 角色服务实现
type RoleServiceImpl struct {
	roleDao                                dao.RoleDao
	userRoleDao                            dao.UserRoleDao
	roleResourceDao                        dao.RoleResourceDao
	roleMenuDao                            dao.RoleMenuDao
	db                                     *gorm.DB
	filterInvocationSecurityMetadataSource *handler.FilterInvocationSecurityMetadataSourceImpl
}

// NewRoleService 创建新的 RoleService 实例
func NewRoleServiceImpl(roleDao dao.RoleDao, userRoleDao dao.UserRoleDao, roleResourceDao dao.RoleResourceDao, roleMenuDao dao.RoleMenuDao, db *gorm.DB, filterInvocationSecurityMetadataSource *handler.FilterInvocationSecurityMetadataSourceImpl) service.RoleService {
	return &RoleServiceImpl{roleDao: roleDao, userRoleDao: userRoleDao, roleResourceDao: roleResourceDao, roleMenuDao: roleMenuDao, db: db, filterInvocationSecurityMetadataSource: filterInvocationSecurityMetadataSource}
}

// ListUserRoles 查询用户角色列表
func (service *RoleServiceImpl) ListUserRoles(ctx context.Context) ([]dto.UserRoleDTO, error) {
	// 查询角色列表
	roleList, err := service.roleDao.ListRoleList(ctx)
	if err != nil {
		return nil, err
	}

	// 将 Role 转换为 UserRoleDTO
	var userRoleDTOs []dto.UserRoleDTO // 创建一个目标类型的切片变量
	userRoleDTOList := utils.BeanCopyList(roleList, userRoleDTOs).([]dto.UserRoleDTO)
	return userRoleDTOList, nil
}

// ListRoles 查询角色列表
func (service *RoleServiceImpl) ListRoles(ctx context.Context, conditionVO *vo.ConditionVO) (*vo.PageResult[dto.RoleDTO], error) {
	// 获取分页参数
	limit := utils.GetLimitCurrent(ctx)
	size := utils.GetSize(ctx)

	// 查询角色列表
	roleDTOList, err := service.roleDao.ListRoles(ctx, limit, size, conditionVO)
	if err != nil {
		return nil, err
	}

	// 查询总量
	count, err := service.roleDao.CountRoles(ctx, conditionVO)
	if err != nil {
		return nil, err
	}

	pageResult := vo.NewPageResult(roleDTOList, int(count))
	return &pageResult, nil
}

// DeleteRoles 删除角色
func (service *RoleServiceImpl) DeleteRoles(ctx context.Context, roleIdList []int) error {
	// 判断角色下是否有用户
	count, err := service.userRoleDao.CountRolesWithUsers(ctx, roleIdList) // 你需要实现 CountByRoleIds 方法
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该角色下存在用户")
	}
	return service.roleDao.DeleteRoles(ctx, roleIdList)
}

// 在 rabbitService/role_service_impl.go 中
func (s *RoleServiceImpl) SaveOrUpdateRole(ctx context.Context, roleVO *vo.RoleVO) error {
	// 检查角色名是否重复
	existRole, err := s.roleDao.CheckRoleExists(ctx, roleVO.RoleName)
	if err != nil {
		return err
	}
	if existRole != nil && existRole.ID != roleVO.ID {
		return errors.New("角色名已存在")
	}

	// 保存或更新角色信息
	role := model.Role{
		ID:        roleVO.ID,
		RoleName:  roleVO.RoleName,
		RoleLabel: roleVO.RoleLabel,
		IsDisable: constants.False, // CommonConst.FALSE
	}
	if err := s.roleDao.SaveOrUpdate(ctx, &role, s.db); err != nil {
		return err
	}

	// 更新角色资源关系
	if roleVO.ResourceIDList != nil {
		if roleVO.ID != 0 {
			if err := s.roleResourceDao.RemoveByRoleID(ctx, roleVO.ID); err != nil {
				return err
			}
		}
		roleResources := make([]model.RoleResource, len(roleVO.ResourceIDList))
		for i, resourceID := range roleVO.ResourceIDList {
			roleResources[i] = model.RoleResource{
				RoleID:     role.ID,
				ResourceID: resourceID,
			}
		}
		if err := s.roleResourceDao.SaveBatch(ctx, roleResources); err != nil {
			return err
		}
		s.filterInvocationSecurityMetadataSource.ClearDataSource() // 重新加载角色资源信息
	}

	// 更新角色菜单关系
	if roleVO.MenuIDList != nil {
		if roleVO.ID != 0 {
			if err := s.roleMenuDao.RemoveByRoleID(ctx, roleVO.ID); err != nil {
				return err
			}
		}
		roleMenus := make([]model.RoleMenu, len(roleVO.MenuIDList))
		for i, menuID := range roleVO.MenuIDList {
			roleMenus[i] = model.RoleMenu{
				RoleID: role.ID,
				MenuID: menuID,
			}
		}
		if err := s.roleMenuDao.SaveBatch(ctx, roleMenus); err != nil {
			return err
		}
	}

	return nil
}
