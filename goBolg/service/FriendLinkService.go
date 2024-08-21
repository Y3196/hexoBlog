package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

type FriendLinkService interface {
	ListFriendLinks(ctx context.Context) ([]dto.FriendLinkDTO, error)

	ListFriendLinkDTO(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.FriendLinkBackDTO], error)

	SaveOrUpdateFriendLink(ctx context.Context, friendLinkVO vo.FriendLinkVO) error

	RemoveFriendLinks(ctx context.Context, ids []uint) error
}
