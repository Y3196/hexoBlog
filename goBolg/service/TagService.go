package service

import (
	"context"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/vo"
)

// TagService 标签服务接口
type TagService interface {
	CountTags(ctx context.Context) (int64, error)

	ListTagsByNames(ctx context.Context, tagNameList []string) ([]model.Tag, error)

	SaveBatch(ctx context.Context, tags []model.Tag) error
	UpdateBatch(ctx context.Context, tags []model.Tag) error
	ListTags(ctx context.Context) (vo.PageResult[dto.TagDTO], error)

	ListTagBackDTO(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.TagBackDTO], error)

	DeleteTag(ctx context.Context, tagIdList []int) error

	SaveOrUpdateTag(ctx context.Context, tagVO vo.TagVO) (model.Tag, error)

	ListTagsBySearch(ctx context.Context, condition vo.ConditionVO) ([]dto.TagDTO, error)
}
