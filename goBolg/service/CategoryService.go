package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

// CategoryService 定义分类服务接口
type CategoryService interface {
	// ListCategories 查询分类列表
	ListCategories(ctx context.Context) vo.PageResult[dto.CategoryDTO]

	// ListBackCategories 查询后台分类列表
	ListBackCategories(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.CategoryBackDTO], error)

	ListCategoriesBySearch(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.CategoryOptionDTO], error)

	SaveOrUpdateCategory(ctx context.Context, categoryVO vo.CategoryVO) error

	DeleteCategory(ctx context.Context, categoryIDList []int) error
}
