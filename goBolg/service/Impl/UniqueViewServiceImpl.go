package Impl

import (
	"context"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/service"
	"time"
)

// uniqueViewServiceImpl 是 UniqueViewService 的具体实现。
type uniqueViewServiceImpl struct {
	uniqueViewDao dao.UniqueViewDao
	redisService  service.RedisService
}

// NewUniqueViewService 创建一个新的 UniqueViewService 实例。
func NewUniqueViewService(uniqueViewDao dao.UniqueViewDao, redisService service.RedisService) service.UniqueViewService {
	return &uniqueViewServiceImpl{
		uniqueViewDao: uniqueViewDao,
		redisService:  redisService,
	}
}

// ListUniqueViews 实现了 UniqueViewService 接口中的 ListUniqueViews 方法。
func (s *uniqueViewServiceImpl) ListUniqueViews(ctx context.Context) ([]dto.UniqueViewDTO, error) {
	startTime := time.Now().AddDate(0, 0, -7) // 7天前
	endTime := time.Now()                     // 现在
	return s.uniqueViewDao.ListUniqueViews(ctx, startTime, endTime)
}

// 下面的定时任务逻辑需要根据你的应用环境来决定是否实现，Go 通常使用第三方库如 cron 进行定时任务。
