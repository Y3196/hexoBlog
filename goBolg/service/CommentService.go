package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

type CommentService interface {
	ListComments(ctx context.Context, commentVO vo.CommentVO, current int, size int) ([]dto.CommentDTO, int, error)

	ListRepliesByCommentId(ctx context.Context, commentId int) ([]dto.ReplyDTO, error)

	SaveComment(ctx context.Context, commentVO vo.CommentVO) (dto.CommentDTO, error)

	// 点赞评论 没测试出来
	SaveCommentLike(ctx context.Context, commentId int) (int, error)

	// 审核评论也没测出来
	UpdateCommentsReview(ctx context.Context, reviewVO vo.ReviewVO) error

	ListCommentBackDTO(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.CommentBackDTO], error)

	RemoveComments(ctx context.Context, commentIdList []int) error
}
