package dao

import (
	"context"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
)

type TagDao interface {
	CountTags(ctx context.Context) (int64, error)

	ListTags(ctx context.Context) ([]model.Tag, error)

	// 查询后台标签列表
	ListTagBackDTO(ctx context.Context, current, size int, condition vo.ConditionVO) ([]dto.TagBackDTO, error)

	// 根据文章id查询标签名
	ListTagNameByArticleId(ctx context.Context, articleId int) ([]string, error)

	// 删除标签
	DeleteTags(ctx context.Context, tagIdList []uint) error

	// 判断标签名是否存在
	CheckTagExists(ctx context.Context, tagName string) (*model.Tag, error)

	GetTagNameByID(ctx context.Context, id int) (string, error)

	ListTagsByNames(ctx context.Context, tagNameList []string) ([]model.Tag, error)

	SaveBatch(ctx context.Context, tags []model.Tag) error

	DeleteBatchIds(ctx context.Context, tagIdList []int) error

	WithContext(ctx context.Context) *gorm.DB
}

type tagDao struct {
	db *gorm.DB
}

func NewTagDao(db *gorm.DB) *tagDao {
	return &tagDao{db}
}

func (dao *tagDao) CreateTag(tag *model.Tag) error {
	return dao.db.Create(tag).Error
}

func (dao *tagDao) GetTagByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := dao.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (dao *tagDao) UpdateTag(tag *model.Tag) error {
	return dao.db.Save(tag).Error
}

func (dao *tagDao) DeleteTag(id uint) error {
	return dao.db.Delete(&model.Tag{}, id).Error
}

// CountTags 实现 CountTags 方法
func (dao *tagDao) CountTags(ctx context.Context) (int64, error) {
	var count int64
	result := dao.db.Model(&model.Tag{}).Count(&count)
	return count, result.Error
}

// ListTags 查询数据库中所有的标签
func (dao *tagDao) ListTags(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	if err := dao.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// ListTagBackDTO 查询后台标签列表
func (dao *tagDao) ListTagBackDTO(ctx context.Context, current, size int, condition vo.ConditionVO) ([]dto.TagBackDTO, error) {
	var tags []dto.TagBackDTO

	// Start building the query
	query := dao.db.WithContext(ctx).
		Table("tb_tag t").
		Select("t.id, t.tag_name, COUNT(tat.article_id) AS article_count, t.create_time").
		Joins("LEFT JOIN tb_article_tag tat ON t.id = tat.tag_id").
		Group("t.id, t.tag_name, t.create_time").
		Order("t.id DESC").
		Offset(current).
		Limit(size)

	// Apply filtering conditions if present
	if condition.Keywords != nil && *condition.Keywords != "" {
		query = query.Where("t.tag_name LIKE ?", "%"+*condition.Keywords+"%")
	}

	if err := query.Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

// ListTagNameByArticleId 根据文章id查询标签名
func (dao *tagDao) ListTagNameByArticleId(ctx context.Context, articleId int) ([]string, error) {
	var tagNames []string
	err := dao.db.WithContext(ctx).
		Raw(`
			SELECT
				tag_name
			FROM
				tb_tag t
				JOIN tb_article_tag tat ON t.id = tat.tag_id
			WHERE
				article_id = ?
		`, articleId).
		Scan(&tagNames).Error

	if err != nil {
		return nil, err
	}
	return tagNames, nil
}

// DeleteTags 删除标签
func (dao *tagDao) DeleteTags(ctx context.Context, tagIdList []uint) error {
	return dao.db.WithContext(ctx).Where("id IN ?", tagIdList).Delete(&model.Tag{}).Error
}

// CheckTagExists 判断标签名是否存在
func (dao *tagDao) CheckTagExists(ctx context.Context, tagName string) (*model.Tag, error) {
	var tag model.Tag
	err := dao.db.WithContext(ctx).
		Where("tag_name = ?", tagName).
		First(&tag).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &tag, nil
}

// getTagNameByID 根据ID获取标签名称
func (dao *tagDao) GetTagNameByID(ctx context.Context, tagID int) (string, error) {
	var name string
	err := dao.db.Table("tb_tag").Select("tag_name").Where("id = ?", tagID).Scan(&name).Error
	if err != nil {
		return "", err
	}
	return name, nil
}

// ListTagsByNames 根据标签名列表查询标签
func (dao *tagDao) ListTagsByNames(ctx context.Context, tagNameList []string) ([]model.Tag, error) {
	var tags []model.Tag
	err := dao.db.WithContext(ctx).
		Where("tag_name IN ?", tagNameList).
		Find(&tags).Error

	if err != nil {
		return nil, err
	}
	return tags, nil
}

// SaveBatch 批量保存标签
func (dao *tagDao) SaveBatch(ctx context.Context, tags []model.Tag) error {
	return dao.db.WithContext(ctx).Create(&tags).Error
}

func (dao *tagDao) WithContext(ctx context.Context) *gorm.DB {
	return dao.db.WithContext(ctx)
}
