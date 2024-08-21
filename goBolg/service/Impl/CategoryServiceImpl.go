package Impl

import (
	"context"
	"errors"
	"fmt"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"time"
)

// categoryServiceImpl 实现了 CategoryService 接口
type categoryServiceImpl struct {
	categoryDao dao.CategoryDao
	articleDao  dao.ArticleDao
}

// NewCategoryService 创建新的 CategoryService 实例
func NewCategoryService(categoryDao dao.CategoryDao, articleDao dao.ArticleDao) service.CategoryService {
	return &categoryServiceImpl{
		categoryDao: categoryDao,
		articleDao:  articleDao,
	}
}

// ListCategories 查询分类列表
// @Summary 查询分类列表
// @Description 获取所有分类及其文章数量
// @Tags Category
// @Produce json
// @Success 200 {object} vo.PageResult[dto.CategoryDTO]
// @Failure 500 {object} vo.Response{error=string}
// @Router /categories [get]
func (s *categoryServiceImpl) ListCategories(ctx context.Context) vo.PageResult[dto.CategoryDTO] {
	categories, err := s.categoryDao.ListCategoryDTO(ctx)
	if err != nil {
		// 处理错误，返回一个空的结果或错误信息
		return vo.PageResult[dto.CategoryDTO]{RecordList: nil, Count: 0}
	}
	count, err := s.categoryDao.CountCategories(ctx)
	if err != nil {
		// 处理错误，返回一个空的结果或错误信息
		return vo.PageResult[dto.CategoryDTO]{RecordList: nil, Count: 0}
	}
	return vo.NewPageResult(categories, int(count))
}

// ListBackCategories 查询后台分类列表
// Debugging in ListBackCategories
func (s *categoryServiceImpl) ListBackCategories(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.CategoryBackDTO], error) {
	count, err := s.categoryDao.CountCategoriesWithCondition(ctx, condition)
	if err != nil {
		return vo.PageResult[dto.CategoryBackDTO]{}, err
	}
	fmt.Printf("Total Count: %d\n", count)

	if count == 0 {
		return vo.NewPageResult([]dto.CategoryBackDTO{}, 0), nil
	}

	limitCurrent := utils.GetLimitCurrent(ctx)
	pageSize := utils.GetSize(ctx)
	categoryList, err := s.categoryDao.ListCategoryBackDTO(ctx, limitCurrent, pageSize, condition)
	if err != nil {
		return vo.PageResult[dto.CategoryBackDTO]{}, err
	}

	fmt.Printf("Category List: %v\n", categoryList)

	return vo.NewPageResult(categoryList, int(count)), nil
}

// ListCategoriesBySearch 根据搜索条件查询分类
// @Summary 根据搜索条件查询分类
// @Description 根据搜索条件获取分类列表
// @Tags Category
// @Produce json
// @Success 200 {object} vo.PageResult[dto.CategoryOptionDTO]
// @Failure 500 {object} vo.Response{error=string}
// @Router /categories/search [get]
func (s *categoryServiceImpl) ListCategoriesBySearch(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.CategoryOptionDTO], error) {
	categories, err := s.categoryDao.ListCategoriesBySearch(ctx, condition)
	if err != nil {
		return vo.PageResult[dto.CategoryOptionDTO]{}, err
	}

	// 假设存在 BeanCopyUtils，在 Go 中我们手动映射
	return vo.NewPageResult(categories, len(categories)), nil
}

// SaveOrUpdateCategory 保存或更新分类
func (s *categoryServiceImpl) SaveOrUpdateCategory(ctx context.Context, categoryVO vo.CategoryVO) error {
	// 检查分类名是否重复
	existingCategory, err := s.categoryDao.GetCategoryByName(ctx, categoryVO.CategoryName)
	if err != nil {
		return err
	}

	// 如果分类名存在且 ID 不匹配，则返回错误
	if existingCategory != nil && existingCategory.ID != categoryVO.ID {
		return errors.New("分类名已存在")
	}

	now := time.Now() // 获取当前时间
	category := &model.Category{
		ID:         categoryVO.ID,
		Name:       categoryVO.CategoryName,
		UpdateTime: now, // 更新操作时设置更新时间
	}

	if existingCategory == nil && category.ID == 0 {
		// 分类名不存在且 ID 为 0，创建新记录
		category.CreateTime = now // 创建时设置创建时间
	}

	// 保存或更新分类
	return s.categoryDao.SaveOrUpdateCategory(ctx, category)
}

// DeleteCategory 删除分类
func (s *categoryServiceImpl) DeleteCategory(ctx context.Context, categoryIDList []int) error {
	// 查询分类ID下是否有文章
	count, err := s.articleDao.GetArticleCountByCategoryIDs(ctx, categoryIDList)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("删除失败，该分类下存在文章")
	}

	// 批量删除分类
	return s.categoryDao.DeleteCategoriesByIDs(ctx, categoryIDList)
}
