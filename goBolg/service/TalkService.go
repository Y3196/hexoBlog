package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

// TalkService 提供说说相关的服务接口
type TalkService interface {
	ListHomeTalks(ctx context.Context) ([]string, error)

	GetTalkById(ctx context.Context, talkId int) (*dto.TalkDTO, error)

	ListTalks(ctx context.Context) (vo.PageResult[dto.TalkDTO], error)

	SaveTalkLike(ctx context.Context, talkId int) (int, error)
}
