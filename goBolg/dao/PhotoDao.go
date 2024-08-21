package dao

import (
	"context"
	constants "goBolg/constant"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
	"time"
)

// PhotoDao 照片 DAO 接口
type PhotoDao interface {
	// 查询照片数量（包括已删除的）
	CountPhotos(ctx context.Context, albumID int) (int64, error)

	// 查询照片数量（不包括已删除的）
	CountActivePhotos(ctx context.Context, albumID uint) (int64, error)

	// 逻辑删除照片
	LogicalDeletePhotos(ctx context.Context, albumID int) error

	// 查询照片列表
	ListPhotos(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]model.Photo, int, error)

	// 根据 ID 更新照片
	UpdateByID(ctx context.Context, photo *model.Photo) error

	//保存照片
	SaveBatch(ctx context.Context, photos []model.Photo) error

	//移动照片相册
	UpdateBatch(ctx context.Context, photos []model.Photo) error

	UpdateBatchByID(ctx context.Context, photos []model.Photo) error

	SelectListByIDs(ctx context.Context, ids []int) ([]model.Photo, error)

	// 删除照片
	DeleteBatchByIDs(ctx context.Context, ids []int) error

	SelectPhotosByAlbumID(ctx context.Context, albumID int, current int, size int) ([]model.Photo, error)
}

type photoDao struct {
	db *gorm.DB
}

// NewPhotoDao 创建新的 PhotoDao 实例
func NewPhotoDao(db *gorm.DB) PhotoDao {
	return &photoDao{db: db}
}

// CountPhotos 查询照片数量（包括已删除的）
func (dao *photoDao) CountPhotos(ctx context.Context, albumID int) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.Photo{}).Where("album_id = ?", albumID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountActivePhotos 查询照片数量（不包括已删除的）
func (dao *photoDao) CountActivePhotos(ctx context.Context, albumID uint) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.Photo{}).Where("album_id = ? AND is_delete = ?", albumID, false).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// LogicalDeletePhotos 逻辑删除照片
func (dao *photoDao) LogicalDeletePhotos(ctx context.Context, albumID int) error {
	return dao.db.WithContext(ctx).Model(&model.Photo{}).Where("album_id = ?", albumID).Update("is_delete", true).Error
}

// ListPhotos 查询照片列表
func (dao *photoDao) ListPhotos(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]model.Photo, int, error) {
	var photos []model.Photo
	query := dao.db.WithContext(ctx).Model(&model.Photo{})

	if condition.AlbumID != nil {
		query = query.Where("album_id = ?", *condition.AlbumID)
	}
	if condition.IsDelete != nil {
		query = query.Where("is_delete = ?", *condition.IsDelete)
	}

	err := query.Order("id DESC, update_time DESC").
		Offset((current - 1) * size).
		Limit(size).
		Find(&photos).
		Error

	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return photos, int(total), nil
}

// UpdateByID 根据 ID 更新照片
func (dao *photoDao) UpdateByID(ctx context.Context, photo *model.Photo) error {
	now := time.Now()
	photo.UpdateTime = &now
	return dao.db.WithContext(ctx).Model(photo).Omit("create_time").Updates(photo).Error
}

func (dao *photoDao) SaveBatch(ctx context.Context, photos []model.Photo) error {
	return dao.db.WithContext(ctx).Create(&photos).Error
}

func (dao *photoDao) UpdateBatch(ctx context.Context, photos []model.Photo) error {
	for _, photo := range photos {
		if err := dao.db.WithContext(ctx).Model(&photo).Updates(map[string]interface{}{
			"album_id": photo.AlbumID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (dao *photoDao) UpdateBatchByID(ctx context.Context, photos []model.Photo) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		for _, photo := range photos {
			if err := tx.Model(&photo).Updates(map[string]interface{}{
				"is_delete":   photo.IsDelete,
				"update_time": now,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (dao *photoDao) SelectListByIDs(ctx context.Context, ids []int) ([]model.Photo, error) {
	var photos []model.Photo
	if err := dao.db.WithContext(ctx).Where("id IN ?", ids).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (dao *photoDao) DeleteBatchByIDs(ctx context.Context, ids []int) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&model.Photo{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (dao *photoDao) SelectPhotosByAlbumID(ctx context.Context, albumID int, current int, size int) ([]model.Photo, error) {
	var photos []model.Photo
	err := dao.db.WithContext(ctx).
		Where("album_id = ? AND is_delete = ?", albumID, constants.False).
		Order("id DESC").
		Offset((current - 1) * size).
		Limit(size).
		Find(&photos).Error
	if err != nil {
		return nil, err
	}
	return photos, nil
}
