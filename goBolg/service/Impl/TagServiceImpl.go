package Impl

import (
	"context"
	"errors"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/utils"
	"goBolg/vo"
	"gorm.io/gorm"
	"time"
)

type tagServiceImpl struct {
	tagDao        dao.TagDao
	articleTagDao dao.ArticleTagDao
	db            *gorm.DB
}

func NewTagServiceImpl(tagDao dao.TagDao, articleTagDao dao.ArticleTagDao, db *gorm.DB) *tagServiceImpl {
	return &tagServiceImpl{tagDao, articleTagDao, db}
}

func (s *tagServiceImpl) CountTags(ctx context.Context) (int64, error) {
	return s.tagDao.CountTags(ctx)
}

func (s *tagServiceImpl) ListTagsByNames(ctx context.Context, tagNameList []string) ([]model.Tag, error) {
	return s.tagDao.ListTagsByNames(ctx, tagNameList)
}

func (s *tagServiceImpl) SaveBatch(ctx context.Context, tags []model.Tag) error {
	return s.tagDao.SaveBatch(ctx, tags)
}

// UpdateBatch 批量更新标签
func (s *tagServiceImpl) UpdateBatch(ctx context.Context, tags []model.Tag) error {
	for _, tag := range tags {
		err := s.tagDao.UpdateTag(ctx, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListTags 查询标签列表
func (s *tagServiceImpl) ListTags(ctx context.Context) (vo.PageResult[dto.TagDTO], error) {
	// 查询标签列表
	tagList, err := s.tagDao.ListTags(ctx)
	if err != nil {
		return vo.PageResult[dto.TagDTO]{}, err
	}

	// 转换 DTO
	tagDTOList := make([]dto.TagDTO, len(tagList))
	for i, tag := range tagList {
		tagDTOList[i] = dto.TagDTO{
			ID:      tag.ID,
			TagName: tag.TagName,
		}
	}

	// 查询标签数量
	count, err := s.tagDao.CountTags(ctx)
	if err != nil {
		return vo.PageResult[dto.TagDTO]{}, err
	}

	return vo.NewPageResult(tagDTOList, int(count)), nil
}

func (s *tagServiceImpl) ListTagBackDTO(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.TagBackDTO], error) {
	// 查询标签数量
	count, err := s.tagDao.CountTags(ctx)
	if err != nil {
		return vo.PageResult[dto.TagBackDTO]{}, err
	}
	if count == 0 {
		return vo.NewPageResult[dto.TagBackDTO](nil, 0), nil
	}

	// 分页查询标签列表
	tagList, err := s.tagDao.ListTagBackDTO(ctx, utils.GetLimitCurrent(ctx), utils.GetSize(ctx), condition)
	if err != nil {
		return vo.PageResult[dto.TagBackDTO]{}, err
	}
	return vo.NewPageResult[dto.TagBackDTO](tagList, int(count)), nil
}

func (s *tagServiceImpl) DeleteTag(ctx context.Context, tagIdList []int) error {
	// 查询标签下是否有文章
	count, err := s.articleTagDao.CountByTagIds(ctx, tagIdList)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("删除失败，该标签下存在文章")
	}

	// 删除标签
	return s.tagDao.DeleteBatchIds(ctx, tagIdList)
}

func (s *tagServiceImpl) SaveOrUpdateTag(ctx context.Context, tagVO vo.TagVO) (model.Tag, error) {
	var tag model.Tag

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 查询标签名是否存在
		var existTag model.Tag
		err := tx.WithContext(ctx).Select("id").Where("tag_name = ?", tagVO.TagName).First(&existTag).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existTag.ID != 0 && existTag.ID != tagVO.ID {
			return errors.New("标签名已存在")
		}

		tag = model.Tag{
			ID:      tagVO.ID,
			TagName: tagVO.TagName,
		}

		if tag.ID == 0 {
			// 新建标签时，设置 create_time 为当前时间
			tag.CreateTime = time.Now()
		} else {
			// 更新标签时，保留原有的 create_time，只更新 update_time
			existingTag := model.Tag{}
			if err := tx.WithContext(ctx).Where("id = ?", tag.ID).First(&existingTag).Error; err != nil {
				return err
			}
			tag.CreateTime = existingTag.CreateTime
			tag.UpdateTime = time.Now()
		}

		if err := tx.WithContext(ctx).Save(&tag).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return model.Tag{}, err
	}

	return tag, nil
}

func (s *tagServiceImpl) ListTagsBySearch(ctx context.Context, condition vo.ConditionVO) ([]dto.TagDTO, error) {
	var tags []model.Tag

	query := s.tagDao.WithContext(ctx).Model(&model.Tag{})

	if condition.Keywords != nil && *condition.Keywords != "" {
		query = query.Where("tag_name LIKE ?", "%"+*condition.Keywords+"%")
	}

	err := query.Order("id DESC").Find(&tags).Error
	if err != nil {
		return nil, err
	}

	// Convert tags to TagDTO
	var tagDTOs []dto.TagDTO
	for _, tag := range tags {
		tagDTOs = append(tagDTOs, dto.TagDTO{
			ID:      tag.ID,
			TagName: tag.TagName,
			// Map other fields if necessary
		})
	}

	return tagDTOs, nil
}
