package dao

import (
	"context"
	"goBolg/dto"
	"goBolg/model"
	"gorm.io/gorm"
	"time"
)

// UniqueViewDao 定义了操作 unique views 数据的接口
type UniqueViewDao interface {
	ListUniqueViews(ctx context.Context, startTime, endTime time.Time) ([]dto.UniqueViewDTO, error)

	InsertUniqueView(ctx context.Context, uniqueView *model.UniqueView) error
}

// uniqueViewDao 实现 UniqueViewDao 接口
type uniqueViewDao struct {
	db *gorm.DB
}

// NewUniqueViewDao 创建一个新的 UniqueViewDao 实例
func NewUniqueViewDao(db *gorm.DB) UniqueViewDao {
	return &uniqueViewDao{db: db}
}

// ListUniqueViews 获取指定时间段内的用户访问数据
func (dao *uniqueViewDao) ListUniqueViews(ctx context.Context, startTime, endTime time.Time) ([]dto.UniqueViewDTO, error) {
	var uniqueViews []dto.UniqueViewDTO
	err := dao.db.Table("tb_unique_view").
		Select("DATE_FORMAT(create_time, '%Y-%m-%d') as day, views_count").
		Where("create_time > ? AND create_time <= ?", startTime, endTime).
		Order("create_time").
		Scan(&uniqueViews).Error
	if err != nil {
		return nil, err
	}
	return uniqueViews, nil
}

// InsertUniqueView 插入新的 unique view 记录
func (dao *uniqueViewDao) InsertUniqueView(ctx context.Context, uniqueView *model.UniqueView) error {
	return dao.db.WithContext(ctx).Create(uniqueView).Error
}
