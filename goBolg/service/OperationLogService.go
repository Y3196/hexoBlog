package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

type OperationLogService interface {
	ListOperationLogs(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.OperationLogDTO], error)
	RemoveOperationLogs(ctx context.Context, ids []uint) error
}
