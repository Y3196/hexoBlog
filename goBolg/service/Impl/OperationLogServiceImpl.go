package Impl

import (
	"context"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"log"
)

type operationLogServiceImpl struct {
	operationLogDao dao.OperationLogDao
}

func NewOperationLogServiceImpl(operationLogDao dao.OperationLogDao) service.OperationLogService {
	return &operationLogServiceImpl{operationLogDao: operationLogDao}
}

func (s *operationLogServiceImpl) ListOperationLogs(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.OperationLogDTO], error) {
	current := utils.GetCurrent(ctx)
	size := utils.GetSize(ctx)
	logs, count, err := s.operationLogDao.Page(ctx, current, size, condition)
	if err != nil {
		log.Printf("Error retrieving operation logs: %v", err)
		return vo.PageResult[dto.OperationLogDTO]{}, err
	}

	logDTOs := utils.BeanCopyList(logs, []dto.OperationLogDTO{}).([]dto.OperationLogDTO)
	return vo.NewPageResult(logDTOs, int(count)), nil
}

func (s *operationLogServiceImpl) RemoveOperationLogs(ctx context.Context, ids []uint) error {
	return s.operationLogDao.RemoveByIds(ctx, ids)
}
