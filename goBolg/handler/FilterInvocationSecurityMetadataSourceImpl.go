package handler

import (
	"context"
	"goBolg/dao"
	"goBolg/dto"
	"log"
	"sync"
)

type FilterInvocationSecurityMetadataSourceImpl struct {
	roleDao dao.RoleDao
}

var resourceRoleList []dto.ResourceRoleDTO
var mu sync.Mutex

func NewFilterInvocationSecurityMetadataSourceImpl(roleDao dao.RoleDao) *FilterInvocationSecurityMetadataSourceImpl {
	return &FilterInvocationSecurityMetadataSourceImpl{
		roleDao: roleDao,
	}
}

// LoadDataSource 加载资源角色信息
func (f *FilterInvocationSecurityMetadataSourceImpl) LoadDataSource(ctx context.Context) {
	mu.Lock()
	defer mu.Unlock()

	resourceRoles, err := f.roleDao.ListResourceRoles(ctx)
	if err != nil {
		log.Printf("加载资源角色信息错误: %v", err)
		return
	}
	resourceRoleList = resourceRoles
}

// ClearDataSource 清空资源角色信息
func (f *FilterInvocationSecurityMetadataSourceImpl) ClearDataSource() {
	mu.Lock()
	defer mu.Unlock()

	resourceRoleList = nil
}

// GetAttributes 获取请求的属性
func (f *FilterInvocationSecurityMetadataSourceImpl) GetAttributes(ctx context.Context, url, method string) []string {
	mu.Lock()
	defer mu.Unlock()

	if len(resourceRoleList) == 0 {
		f.LoadDataSource(ctx)
	}

	for _, resourceRole := range resourceRoleList {
		if resourceRole.URL == url && resourceRole.RequestMethod == method {
			return resourceRole.RoleList
		}
	}

	return nil
}
