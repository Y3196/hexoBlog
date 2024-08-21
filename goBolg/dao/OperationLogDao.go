package dao

import (
	"context"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
)

// OperationLogDao 操作日志 DAO 接口
type OperationLogDao interface {
	Page(ctx context.Context, page int, size int, condition vo.ConditionVO) ([]model.OperationLog, int64, error)

	RemoveByIds(ctx context.Context, ids []uint) error
	// 插入一条操作日志记录
	//Insert(ctx contxt.Context, operationLog model.OperationLog) error
}

type operationLogDao struct {
	db *gorm.DB
}

// NewOperationLogDao 创建新的 OperationLogDao 实例
func NewOperationLogDao(db *gorm.DB) OperationLogDao {
	return &operationLogDao{db: db}
}

func (dao *operationLogDao) Page(ctx context.Context, page int, size int, condition vo.ConditionVO) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var count int64

	offset := (page - 1) * size
	query := dao.db.WithContext(ctx).Model(&model.OperationLog{}).Offset(offset).Limit(size)

	if condition.Keywords != nil && *condition.Keywords != "" {
		query = query.Where("opt_module LIKE ?", "%"+*condition.Keywords+"%").
			Or("opt_desc LIKE ?", "%"+*condition.Keywords+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("id DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, count, nil
}

// Insert 插入一条操作日志记录
/*func (dao *operationLogDao) Insert(ctx contxt.Context, operationLog model.OperationLog) error {
	return dao.db.WithContext(ctx).Create(&operationLog).Error
}
*/
func (dao *operationLogDao) RemoveByIds(ctx context.Context, ids []uint) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.OperationLog{}).Error
}
