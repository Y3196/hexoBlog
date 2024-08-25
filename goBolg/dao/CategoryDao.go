package dao

import (
	"context"
	"errors"
	"fmt"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
)

type CategoryDao interface {
	CountCategories(ctx context.Context) (int64, error)
	CountCategoriesWithCondition(ctx context.Context, condition vo.ConditionVO) (int64, error)

	ListCategoryDTO(ctx context.Context) ([]dto.CategoryDTO, error)

	ListCategoryBackDTO(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]dto.CategoryBackDTO, error)

	GetCategoryNameByID(ctx context.Context, id int) (string, error)

	GetCategoryByID(ctx context.Context, categoryID int) (*model.Category, error)

	ListCategoriesBySearch(ctx context.Context, condition vo.ConditionVO) ([]dto.CategoryOptionDTO, error)

	GetCategoryByName(ctx context.Context, categoryName string) (*model.Category, error)

	InsertCategory(ctx context.Context, category *model.Category) error

	ExistsCategoryByName(ctx context.Context, categoryName string) (bool, error)

	SaveOrUpdateCategory(ctx context.Context, category *model.Category) error

	DeleteCategoriesByIDs(ctx context.Context, categoryIDList []int) error
}

type categoryDao struct {
	db *gorm.DB
}

func NewCategoryDao(db *gorm.DB) CategoryDao {
	return &categoryDao{db: db}
}

func (dao *categoryDao) CountCategories(ctx context.Context) (int64, error) {
	var count int64
	err := dao.db.Table("tb_category").Count(&count).Error
	return count, err
}

// BuildCategoryQuery 构建分类查询
func (dao *categoryDao) BuildCategoryQuery(query *gorm.DB, condition vo.ConditionVO) *gorm.DB {
	// 处理 Keywords 指针
	if condition.Keywords != nil && *condition.Keywords != "" {
		query = query.Where("c.category_name LIKE ?", "%"+*condition.Keywords+"%")
	}
	return query
}

func (dao *categoryDao) ListCategoryDTO(ctx context.Context) ([]dto.CategoryDTO, error) {
	var categories []dto.CategoryDTO
	query := `
		SELECT
			c.id,
			c.category_name,
			COUNT(a.id) AS article_count
		FROM
			tb_category c
			LEFT JOIN (SELECT id, category_id FROM tb_article WHERE is_delete = 0 AND status = 1) a ON c.id = a.category_id
		GROUP BY
			c.id`
	err := dao.db.Raw(query).Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (dao *categoryDao) ListCategoryBackDTO(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]dto.CategoryBackDTO, error) {
	var categories []dto.CategoryBackDTO

	offset := (current - 1) * size

	// 构建 SQL 查询
	query := dao.db.Table("tb_category c").
		Select("c.id, c.category_name, COALESCE(COUNT(a.id), 0) AS article_count, c.create_time").
		Joins("LEFT JOIN tb_article a ON c.id = a.category_id AND a.is_delete = 0 AND a.status = 1").
		Group("c.id").
		Order("c.id DESC").
		Offset(offset).
		Limit(size)

	// 使用 BuildCategoryQuery 添加条件
	query = dao.BuildCategoryQuery(query, condition)

	// 添加调试日志
	query = query.Debug()

	err := query.Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	fmt.Printf("Offset: %d, Limit: %d, Categories: %v\n", offset, size, categories)

	return categories, nil
}

func (dao *categoryDao) GetCategoryNameByID(ctx context.Context, categoryID int) (string, error) {
	var name string
	err := dao.db.Table("tb_category").Select("category_name").Where("id = ?", categoryID).Scan(&name).Error
	if err != nil {
		return "", err
	}
	return name, nil
}

// GetCategoryByID 根据 ID 获取分类
func (dao *categoryDao) GetCategoryByID(ctx context.Context, categoryID int) (*model.Category, error) {
	var category model.Category
	err := dao.db.WithContext(ctx).Table("tb_category").Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录未找到，返回 nil 和 nil 错误
			return nil, nil
		}
		// 其他错误
		return nil, err
	}
	return &category, nil
}

// GetCategoryByName 根据名称获取分类
func (dao *categoryDao) GetCategoryByName(ctx context.Context, categoryName string) (*model.Category, error) {
	var category model.Category
	err := dao.db.WithContext(ctx).Table("tb_category").Where("category_name = ?", categoryName).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (dao *categoryDao) InsertCategory(ctx context.Context, category *model.Category) error {
	return dao.db.WithContext(ctx).Create(category).Error
}

func (dao *categoryDao) CountCategoriesWithCondition(ctx context.Context, condition vo.ConditionVO) (int64, error) {
	var count int64
	query := dao.db.Table("tb_category c").Select("COUNT(DISTINCT c.id)")

	// 使用 BuildCategoryQuery 添加条件
	query = dao.BuildCategoryQuery(query, condition)

	// 添加调试日志
	query = query.Debug()

	err := query.Count(&count).Error
	return count, err
}

func (dao *categoryDao) ListCategoriesBySearch(ctx context.Context, condition vo.ConditionVO) ([]dto.CategoryOptionDTO, error) {
	var categories []dto.CategoryOptionDTO

	// 初始化关键词，处理可能的 nil 值
	var keywords string
	if condition.Keywords != nil && *condition.Keywords != "" {
		keywords = "%" + *condition.Keywords + "%"
	} else {
		keywords = "%" // 如果 keywords 为 nil 或者空字符串, 就匹配所有记录
	}

	// 查询构建
	query := dao.db.Table("tb_category").
		Select("id, category_name, create_time").
		Where("category_name LIKE ?", keywords).
		Order("id DESC")

	// 添加调试日志
	query = query.Debug()

	// 执行查询
	err := query.Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (dao *categoryDao) ExistsCategoryByName(ctx context.Context, categoryName string) (bool, error) {
	var count int64
	err := dao.db.Table("tb_category").
		Where("category_name = ?", categoryName).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SaveOrUpdateCategory 插入或更新分类
func (dao *categoryDao) SaveOrUpdateCategory(ctx context.Context, category *model.Category) error {
	if category.ID == 0 {
		// 插入新记录时，忽略更新时间
		return dao.db.WithContext(ctx).Omit("UpdateTime").Create(category).Error
	}
	// 更新记录时，只更新分类名和更新时间，忽略创建时间
	return dao.db.WithContext(ctx).Model(&model.Category{}).
		Where("id = ?", category.ID).
		Updates(map[string]interface{}{
			"category_name": category.Name,
			"update_time":   category.UpdateTime,
		}).Error
}

// DeleteCategoriesByIDs 批量删除分类ID
func (dao *categoryDao) DeleteCategoriesByIDs(ctx context.Context, categoryIDList []int) error {
	return dao.db.WithContext(ctx).Where("id IN ?", categoryIDList).Delete(&model.Category{}).Error
}
