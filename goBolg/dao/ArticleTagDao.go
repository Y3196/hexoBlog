package dao

import (
	"context"
	"goBolg/model"
	"gorm.io/gorm"
)

// ArticleTagDao 文章标签 DAO 接口
type ArticleTagDao interface {
	DeleteByArticleIds(ctx context.Context, articleIdList []int) error

	DeleteByArticleId(ctx context.Context, articleId int) error

	CountByArticleId(ctx context.Context, articleId uint) (int64, error)

	// SaveBatch 批量保存 ArticleTag 实体到数据库
	SaveBatch(articleTags []model.ArticleTag) error

	CountByTagIds(ctx context.Context, tagIdList []int) (int64, error)
}

type articleTagDao struct {
	db *gorm.DB
}

// NewArticleTagDao 创建新的 ArticleTagDao 实例
func NewArticleTagDao(db *gorm.DB) ArticleTagDao {
	return &articleTagDao{db: db}
}

// DeleteByArticleIds 根据文章 ID 列表删除文章标签
func (dao *articleTagDao) DeleteByArticleIds(ctx context.Context, articleIdList []int) error {
	if len(articleIdList) == 0 {
		return nil
	}

	// 执行删除操作
	result := dao.db.WithContext(ctx).Where("article_id IN ?", articleIdList).Delete(&model.ArticleTag{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteByArticleId 根据文章 ID 删除文章标签
func (dao *articleTagDao) DeleteByArticleId(ctx context.Context, articleId int) error {
	// 执行删除操作
	result := dao.db.WithContext(ctx).Where("article_id = ?", articleId).Delete(&model.ArticleTag{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CountByArticleId 根据文章 ID 统计文章标签数量
func (dao *articleTagDao) CountByArticleId(ctx context.Context, articleId uint) (int64, error) {
	var count int64
	result := dao.db.WithContext(ctx).Model(&model.ArticleTag{}).Where("article_id = ?", articleId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// SaveBatch 批量保存 ArticleTag 实体到数据库
func (r *articleTagDao) SaveBatch(articleTags []model.ArticleTag) error {
	return r.db.Create(&articleTags).Error
}

// CountByTagIds 根据标签ID列表统计文章标签数量
func (dao *articleTagDao) CountByTagIds(ctx context.Context, tagIdList []int) (int64, error) {
	var count int64
	result := dao.db.WithContext(ctx).Model(&model.ArticleTag{}).Where("tag_id IN ?", tagIdList).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// DeleteBatchIds 批量删除标签
func (dao *tagDao) DeleteBatchIds(ctx context.Context, tagIdList []int) error {
	return dao.db.WithContext(ctx).Where("id IN ?", tagIdList).Delete(&model.Tag{}).Error
}
