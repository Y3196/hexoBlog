package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

type MessageService interface {
	SaveMessage(ctx context.Context, messageVO vo.MessageVO) error

	ListMessages(ctx context.Context) ([]dto.MessageDTO, error)

	ListMessageBackDTO(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.MessageBackDTO], error)

	UpdateMessagesReview(ctx context.Context, reviewVO vo.ReviewVO) error

	DeleteMessages(ctx context.Context, messageIdList []uint) error
}
