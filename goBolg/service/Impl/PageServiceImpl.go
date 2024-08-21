package Impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"goBolg/dao"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"gorm.io/gorm"
	"time"
)

type pageServiceImpl struct {
	pageDao      dao.PageDao
	redisService service.RedisService
	db           *gorm.DB
}

func NewPageService(db *gorm.DB, pageDao dao.PageDao, redisService *RedisServiceImpl) service.PageService {
	return &pageServiceImpl{
		db:           db,
		pageDao:      pageDao,
		redisService: redisService,
	}
}

func (s *pageServiceImpl) SaveOrUpdatePage(ctx context.Context, pageVO vo.PageVO) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	page := &model.Page{}
	utils.BeanCopyObject(pageVO, page)

	if page.ID != 0 {
		var existingPage model.Page
		err := s.db.WithContext(ctx).First(&existingPage, page.ID).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return err
		}
		if err == nil {
			page.CreateTime = existingPage.CreateTime
			now := time.Now()
			page.UpdateTime = &now
		} else {
			page.CreateTime = time.Now()
		}
	} else {
		page.CreateTime = time.Now()
	}

	if err := tx.WithContext(ctx).Save(page).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	s.redisService.Del(ctx, "page_cover")
	return nil
}

func (s *pageServiceImpl) DeletePage(ctx context.Context, pageID uint) error {
	if err := s.pageDao.Delete(ctx, pageID); err != nil {
		return err
	}
	s.redisService.Del(ctx, "page_cover")
	return nil
}

func (s *pageServiceImpl) ListPages(ctx context.Context) ([]vo.PageVO, error) {
	var pageVOList []vo.PageVO
	pageList, err := s.redisService.Get(ctx, "page_cover")
	if err == redis.Nil {
		// 缓存未命中，从数据库获取
		pages, dbErr := s.pageDao.List(ctx)
		if dbErr != nil {
			return nil, dbErr
		}
		for _, page := range pages {
			pageVO := vo.PageVO{
				ID:        int(page.ID),
				PageName:  page.PageName,
				PageLabel: page.PageLabel,
				PageCover: page.PageCover,
			}
			pageVOList = append(pageVOList, pageVO)
		}
		jsonData, jsonErr := json.Marshal(pageVOList)
		if jsonErr != nil {
			return nil, jsonErr
		}
		if redisErr := s.redisService.Set(ctx, "page_cover", jsonData, 0); redisErr != nil {
			return nil, redisErr
		}
	} else if err != nil {
		return nil, err
	} else {
		// 打印获取的原始数据
		fmt.Println("Raw pageList data from Redis:", pageList)

		// 直接反序列化，不需要双重解码
		if err := json.Unmarshal([]byte(pageList), &pageVOList); err != nil {
			// 提供详细的错误信息
			fmt.Printf("Error unmarshalling pageList data: %v\n", err)
			return nil, err
		}
	}

	return pageVOList, nil
}
