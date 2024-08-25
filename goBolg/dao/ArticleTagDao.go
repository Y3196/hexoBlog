package dao

import (
	"context"
	"goBolg/model"
	"gorm.io/gorm"
	"time"
)

// ArticleTagDao 文章标签 DAO 接口
type ArticleTagDao interface {
	DeleteByArticleIds(ctx context.Context, articleIdList []int) error

	DeleteByArticleId(ctx context.Context, articleId int) error

	CountByArticleId(ctx context.Context, articleId uint) (int64, error)

	// SaveBatch 批量保存 ArticleTag 实体到数据库
	SaveBatch(articleTags []model.ArticleTag) error

	CountByTagIds(ctx context.Context, tagIdList []int) (int64, error)

	ExistsByArticleAndTag(ctx context.Context, articleId int, tagId int) (bool, error)

	ListByArticleId(ctx context.Context, articleId int) ([]model.ArticleTag, error)

	UpdateTimestamp(ctx context.Context, articleId int, tagId int, timestamp time.Time) error
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
func (dao *articleTagDao) ExistsByArticleAndTag(ctx context.Context, articleId int, tagId int) (bool, error) {
	var count int64
	sql := "SELECT COUNT(*) FROM tb_article_tag WHERE article_id = ? AND tag_id = ?"
	err := dao.db.Raw(sql, articleId, tagId).Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *articleTagDao) ListByArticleId(ctx context.Context, articleId int) ([]model.ArticleTag, error) {
	var articleTags []model.ArticleTag
	err := dao.db.WithContext(ctx).
		Where("article_id = ?", articleId).
		Find(&articleTags).Error
	if err != nil {
		return nil, err
	}
	return articleTags, nil
}

// UpdateTimestamp 更新现有的文章-标签关联的更新时间
func (dao *articleTagDao) UpdateTimestamp(ctx context.Context, articleId int, tagId int, updateTime time.Time) error {
	return dao.db.WithContext(ctx).
		Model(&model.ArticleTag{}).
		Where("article_id = ? AND tag_id = ?", articleId, tagId).
		Update("updated_at", updateTime).Error
}
