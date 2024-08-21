package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

type ResourceService interface {
	SaveOrUpdateResource(ctx context.Context, resourceVO vo.ResourceVO) error

	DeleteResource(ctx context.Context, resourceId uint) error

	ListResources(ctx context.Context, conditionVO vo.ConditionVO) ([]dto.ResourceDTO, error)

	ListResourceOption(ctx context.Context) ([]dto.LabelOptionDTO, error)
}
