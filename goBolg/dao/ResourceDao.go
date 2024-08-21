package dao

import (
	"context"
	"goBolg/model"
	"gorm.io/gorm"
	"time"
)

// ResourceDao 资源 DAO 接口
type ResourceDao interface {
	// 删除子资源
	DeleteChildResources(ctx context.Context, parentID uint) error

	// 查询资源列表（带条件）
	ListResources(ctx context.Context, keywords *string, result *[]model.Resource) error

	// 查询资源列表（带选择字段和条件）
	ListResourcesWithFields(ctx context.Context, isAnonymous bool) ([]model.Resource, error)

	// 保存或更新资源
	SaveOrUpdateResource(ctx context.Context, resource *model.Resource) error

	ListChildResourceIds(ctx context.Context, parentID uint) ([]uint, error)

	DeleteResources(ctx context.Context, ids []uint) error
}

type resourceDao struct {
	db *gorm.DB
}

// NewResourceDao 创建新的 ResourceDao 实例
func NewResourceDao(db *gorm.DB) ResourceDao {
	return &resourceDao{db: db}
}

// DeleteChildResources 删除子资源
func (dao *resourceDao) DeleteChildResources(ctx context.Context, parentID uint) error {
	// 查询所有子资源 ID
	var resourceIds []uint
	err := dao.db.WithContext(ctx).
		Model(&model.Resource{}).
		Where("parent_id = ?", parentID).
		Pluck("id", &resourceIds).Error
	if err != nil {
		return err
	}
	// 删除子资源
	return dao.db.WithContext(ctx).Where("id IN ?", resourceIds).Delete(&model.Resource{}).Error
}

// ListResources 查询资源列表（带条件）
func (dao *resourceDao) ListResources(ctx context.Context, keywords *string, result *[]model.Resource) error {
	query := dao.db.WithContext(ctx).Model(&model.Resource{})
	if keywords != nil && *keywords != "" {
		query = query.Where("resource_name LIKE ?", "%"+*keywords+"%")
	}
	return query.Find(result).Error
}

// ListResourcesWithFields 查询资源列表（带选择字段和条件）
func (dao *resourceDao) ListResourcesWithFields(ctx context.Context, isAnonymous bool) ([]model.Resource, error) {
	var resources []model.Resource
	err := dao.db.WithContext(ctx).
		Model(&model.Resource{}).
		Select("id, resource_name, parent_id").
		Where("is_anonymous = ?", isAnonymous).
		Find(&resources).Error
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// SaveOrUpdateResource 保存或更新资源
func (dao *resourceDao) SaveOrUpdateResource(ctx context.Context, resource *model.Resource) error {
	var existingResource model.Resource
	err := dao.db.WithContext(ctx).First(&existingResource, resource.ID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 记录不存在，设置创建时间并保存新记录
			resource.CreateTime = time.Now()
			// 确保 UpdateTime 为 nil
			resource.UpdateTime = nil
			return dao.db.WithContext(ctx).Create(resource).Error
		}
		return err
	}

	// 记录存在，更新字段并设置更新时间
	if resource.ResourceName != "" {
		existingResource.ResourceName = resource.ResourceName
	}
	if resource.URL != "" {
		existingResource.URL = resource.URL
	}
	if resource.RequestMethod != "" {
		existingResource.RequestMethod = resource.RequestMethod
	}
	if *resource.ParentID != 0 {
		existingResource.ParentID = resource.ParentID
	}
	existingResource.IsAnonymous = resource.IsAnonymous

	// 只在更新时设置更新时间
	now := time.Now()
	existingResource.UpdateTime = &now

	return dao.db.WithContext(ctx).Save(&existingResource).Error
}

func (dao *resourceDao) ListChildResourceIds(ctx context.Context, parentID uint) ([]uint, error) {
	var ids []uint
	err := dao.db.WithContext(ctx).
		Model(&model.Resource{}).
		Where("parent_id = ?", parentID).
		Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (dao *resourceDao) DeleteResources(ctx context.Context, ids []uint) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Resource{}).Error
}
