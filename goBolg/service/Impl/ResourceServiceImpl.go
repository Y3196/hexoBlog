package Impl

import (
	"context"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/exception"
	"goBolg/handler"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
)

type ResourceServiceImpl struct {
	resourceDao                            dao.ResourceDao
	roleResourceDao                        dao.RoleResourceDao
	filterInvocationSecurityMetadataSource *handler.FilterInvocationSecurityMetadataSourceImpl
}

func NewResourceService(resourceDao dao.ResourceDao, roleResourceDao dao.RoleResourceDao, filterInvocationSecurityMetadataSource *handler.FilterInvocationSecurityMetadataSourceImpl) service.ResourceService {
	return &ResourceServiceImpl{
		resourceDao:                            resourceDao,
		roleResourceDao:                        roleResourceDao,
		filterInvocationSecurityMetadataSource: filterInvocationSecurityMetadataSource,
	}
}

func (s *ResourceServiceImpl) SaveOrUpdateResource(ctx context.Context, resourceVO vo.ResourceVO) error {
	resource := &model.Resource{}
	utils.BeanCopy(resourceVO, resource)

	if err := s.resourceDao.SaveOrUpdateResource(ctx, resource); err != nil {
		return err
	}

	// 重新加载角色资源信息
	s.filterInvocationSecurityMetadataSource.ClearDataSource()
	return nil
}

func (s *ResourceServiceImpl) DeleteResource(ctx context.Context, resourceId uint) error {
	// 查询是否有角色关联
	count, err := s.roleResourceDao.CountByResourceID(ctx, resourceId)
	if err != nil {
		return err
	}
	if count > 0 {
		return exception.NewBizException("该资源下存在角色")
	}

	// 查询所有子资源ID
	resourceIds, err := s.resourceDao.ListChildResourceIds(ctx, resourceId)
	if err != nil {
		return err
	}
	resourceIds = append(resourceIds, resourceId)

	// 删除子资源及当前资源
	return s.resourceDao.DeleteResources(ctx, resourceIds)
}

func (s *ResourceServiceImpl) ListResources(ctx context.Context, conditionVO vo.ConditionVO) ([]dto.ResourceDTO, error) {
	// 查询资源列表
	var resourceList []model.Resource
	err := s.resourceDao.ListResources(ctx, conditionVO.Keywords, &resourceList)
	if err != nil {
		return nil, err
	}

	// 获取所有模块
	parentList := listResourceModule(resourceList)
	// 根据父ID分组获取模块下的资源
	childrenMap := listResourceChildren(resourceList)

	// 绑定模块下的所有接口
	var resourceDTOList []dto.ResourceDTO
	for _, item := range parentList {
		resourceDTO := dto.ResourceDTO{
			ID:            item.ID,
			ResourceName:  item.ResourceName,
			URL:           item.URL,
			RequestMethod: item.RequestMethod,
			IsDisable:     item.IsAnonymous, // Assuming `IsDisable` corresponds to `IsAnonymous`
			CreateTime:    item.CreateTime,
		}
		if children, ok := childrenMap[item.ID]; ok {
			var childrenDTOList []dto.ResourceDTO
			for _, child := range children {
				childDTO := dto.ResourceDTO{
					ID:            child.ID,
					ResourceName:  child.ResourceName,
					URL:           child.URL,
					RequestMethod: child.RequestMethod,
					IsDisable:     child.IsAnonymous, // Assuming `IsDisable` corresponds to `IsAnonymous`
					CreateTime:    child.CreateTime,
				}
				childrenDTOList = append(childrenDTOList, childDTO)
			}
			resourceDTO.Children = childrenDTOList
		}
		resourceDTOList = append(resourceDTOList, resourceDTO)
		delete(childrenMap, item.ID)
	}

	// 若还有资源未取出则拼接
	if len(childrenMap) > 0 {
		var childrenList []model.Resource
		for _, children := range childrenMap {
			childrenList = append(childrenList, children...)
		}
		for _, child := range childrenList {
			childDTO := dto.ResourceDTO{
				ID:            child.ID,
				ResourceName:  child.ResourceName,
				URL:           child.URL,
				RequestMethod: child.RequestMethod,
				IsDisable:     child.IsAnonymous, // Assuming `IsDisable` corresponds to `IsAnonymous`
				CreateTime:    child.CreateTime,
			}
			resourceDTOList = append(resourceDTOList, childDTO)
		}
	}

	return resourceDTOList, nil
}

// listResourceOption 列出资源选项
func (s *ResourceServiceImpl) ListResourceOption(ctx context.Context) ([]dto.LabelOptionDTO, error) {
	// 查询资源列表
	var resourceList []model.Resource

	// Pass nil for keywords if you want to retrieve all resources
	err := s.resourceDao.ListResources(ctx, nil, &resourceList)
	if err != nil {
		return nil, err
	}

	// 过滤掉匿名资源
	var nonAnonymousResources []model.Resource
	for _, resource := range resourceList {
		if resource.IsAnonymous == 0 { // Assuming 0 means FALSE
			nonAnonymousResources = append(nonAnonymousResources, resource)
		}
	}

	// 获取所有模块
	parentList := listResourceModule(nonAnonymousResources)
	// 根据父ID分组获取模块下的资源
	childrenMap := listResourceChildren(nonAnonymousResources)

	// 组装父子数据
	var optionsList []dto.LabelOptionDTO
	for _, item := range parentList {
		var childrenList []dto.LabelOptionDTO
		if children, ok := childrenMap[item.ID]; ok {
			for _, child := range children {
				childrenList = append(childrenList, dto.LabelOptionDTO{
					ID:    child.ID,
					Label: child.ResourceName,
				})
			}
		}
		optionsList = append(optionsList, dto.LabelOptionDTO{
			ID:       item.ID,
			Label:    item.ResourceName,
			Children: childrenList,
		})
	}

	return optionsList, nil
}

// listResourceModule 获取所有资源模块
func listResourceModule(resourceList []model.Resource) []model.Resource {
	var modules []model.Resource
	for _, item := range resourceList {
		if item.ParentID == nil {
			modules = append(modules, item)
		}
	}
	return modules
}

// listResourceChildren 获取模块下的所有资源
func listResourceChildren(resourceList []model.Resource) map[uint][]model.Resource {
	childrenMap := make(map[uint][]model.Resource)
	for _, item := range resourceList {
		if item.ParentID != nil {
			parentID := *item.ParentID
			childrenMap[parentID] = append(childrenMap[parentID], item)
		}
	}
	return childrenMap
}
