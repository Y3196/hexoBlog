package dao

import (
	"context"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
)

// MessageDao 定义与消息操作相关的方法
type MessageDao interface {
	CountMessages(ctx context.Context) (int64, error)
	// 插入一条新记录
	Insert(ctx context.Context, message model.Message) error

	// 查询留言列表
	ListMessages(ctx context.Context, isReview bool) ([]model.Message, error)

	// 分页查询留言
	PageMessages(ctx context.Context, page, size int, condition vo.ConditionVO) ([]model.Message, int64, error)

	// 更新审核留言
	UpdateBatchById(ctx context.Context, messages []model.Message) error

	// 删除留言
	RemoveByIds(ctx context.Context, ids []uint) error
}

// messageDaoImpl 是 MessageDao 的一个具体实现
type messageDaoImpl struct {
	db *gorm.DB
}

// NewMessageDao 创建并返回一个实现 MessageDao 的新实例
func NewMessageDao(db *gorm.DB) MessageDao {
	return &messageDaoImpl{db: db}
}

// CountMessages 实现 MessageDao 中的方法，返回当前数据库中留言的总数
func (dao *messageDaoImpl) CountMessages(ctx context.Context) (int64, error) {
	var count int64
	result := dao.db.Table("tb_message").Count(&count)
	return count, result.Error
}

// Insert 插入一条新记录
func (dao *messageDaoImpl) Insert(ctx context.Context, message model.Message) error {
	return dao.db.WithContext(ctx).Omit("UpdateTime").Create(&message).Error
}

// ListMessages 查询留言列表
func (dao *messageDaoImpl) ListMessages(ctx context.Context, isReview bool) ([]model.Message, error) {
	var messages []model.Message
	err := dao.db.WithContext(ctx).
		Select("id, nickname, avatar, message_content, time").
		Where("is_review = ?", isReview).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// PageMessages 分页查询留言
func (dao *messageDaoImpl) PageMessages(ctx context.Context, page, size int, condition vo.ConditionVO) ([]model.Message, int64, error) {
	var messages []model.Message
	var count int64

	offset := (page - 1) * size

	query := dao.db.WithContext(ctx).Model(&model.Message{})

	if condition.Keywords != nil && *condition.Keywords != "" {
		query = query.Where("nickname LIKE ?", "%"+*condition.Keywords+"%")
	}

	if condition.IsReview != nil {
		query = query.Where("is_review = ?", *condition.IsReview)
	}

	// 查询总记录数
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询分页记录
	err = query.Order("id DESC").Offset(offset).Limit(size).Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, count, nil
}

// UpdateBatchById 批量更新留言
func (dao *messageDaoImpl) UpdateBatchById(ctx context.Context, messages []model.Message) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, message := range messages {
			if err := tx.Model(&model.Message{}).Where("id = ?", message.ID).Updates(map[string]interface{}{
				"is_review": message.IsReview,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveByIds 删除留言
func (dao *messageDaoImpl) RemoveByIds(ctx context.Context, ids []uint) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Message{}).Error
}
