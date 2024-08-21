package dao

import (
	"context"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
	"time"
)

// PhotoAlbumDao 相册 DAO 接口
type PhotoAlbumDao interface {
	// 查询后台相册列表
	ListPhotoAlbumBacks(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]dto.PhotoAlbumBackDTO, error)

	// 查询相册名是否存在
	FindAlbumByName(ctx context.Context, albumName string) (*model.PhotoAlbum, error)

	// 查询相册数量
	CountPhotoAlbums(ctx context.Context, condition vo.ConditionVO) (int, error)

	// 查询相册列表
	ListPhotoAlbums(ctx context.Context) ([]model.PhotoAlbum, error)

	// 查询相册信息
	GetPhotoAlbumByID(ctx context.Context, albumID int) (*model.PhotoAlbum, error)

	// 逻辑删除相册
	LogicalDeletePhotoAlbum(ctx context.Context, albumID int) error

	// 删除相册
	DeletePhotoAlbum(ctx context.Context, albumID int) error

	// 查询相册列表（发布状态）
	ListPublishedPhotoAlbums(ctx context.Context) ([]model.PhotoAlbum, error)

	// 根据ID更新相册
	UpdatePhotoAlbumByID(ctx context.Context, albumID int, updates map[string]interface{}) error

	// 保存相册
	SavePhotoAlbum(ctx context.Context, album *model.PhotoAlbum) error

	// 更新相册
	UpdatePhotoAlbum(ctx context.Context, album *model.PhotoAlbum) error

	//根据ID查询相册中的数量
	CountPhotosInAlbumById(ctx context.Context, albumID int) (int64, error)

	UpdateBatchByID(ctx context.Context, albums []model.PhotoAlbum) error

	FindAlbumByID(ctx context.Context, albumID int, isDelete int, status int) (*model.PhotoAlbum, error)
}

type photoAlbumDao struct {
	db *gorm.DB
}

// NewPhotoAlbumDao 创建新的 PhotoAlbumDao 实例
func NewPhotoAlbumDao(db *gorm.DB) PhotoAlbumDao {
	return &photoAlbumDao{db: db}
}

// ListPhotoAlbumBacks 查询后台相册列表
func (dao *photoAlbumDao) ListPhotoAlbumBacks(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]dto.PhotoAlbumBackDTO, error) {
	var result []dto.PhotoAlbumBackDTO
	query := `
		SELECT
			pa.id,
			album_name,
			album_desc,
			album_cover,
			COUNT(a.id) AS photo_count,
			status
		FROM
			(
				SELECT
					id,
					album_name,
					album_desc,
					album_cover,
					status
				FROM
					tb_photo_album
				WHERE
					is_delete = 0
				`

	var queryParams []interface{}
	if condition.Keywords != nil && *condition.Keywords != "" {
		query += ` AND album_name LIKE ?`
		queryParams = append(queryParams, "%"+*condition.Keywords+"%")
	}

	query += `
				ORDER BY
					id DESC
				LIMIT ?, ?
			) pa
		LEFT JOIN
			(
				SELECT
					id,
					album_id
				FROM
					tb_photo
				WHERE
					is_delete = 0
			) a ON pa.id = a.album_id
		GROUP BY
			pa.id
	`

	queryParams = append(queryParams, current, size)
	err := dao.db.Raw(query, queryParams...).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindAlbumByName 查询相册名是否存在
func (dao *photoAlbumDao) FindAlbumByName(ctx context.Context, albumName string) (*model.PhotoAlbum, error) {
	var album model.PhotoAlbum
	err := dao.db.WithContext(ctx).Where("album_name = ?", albumName).First(&album).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

// CountPhotoAlbums 查询相册数量
func (dao *photoAlbumDao) CountPhotoAlbums(ctx context.Context, condition vo.ConditionVO) (int, error) {
	var count int64
	query := dao.db.WithContext(ctx).Model(&model.PhotoAlbum{})
	if condition.Keywords != nil && *condition.Keywords != "" {
		query = query.Where("album_name LIKE ?", "%"+*condition.Keywords+"%")
	}
	query = query.Where("is_delete = ?", false).Count(&count)
	return int(count), query.Error
}

// ListPhotoAlbums 查询相册列表
func (dao *photoAlbumDao) ListPhotoAlbums(ctx context.Context) ([]model.PhotoAlbum, error) {
	var albums []model.PhotoAlbum
	err := dao.db.WithContext(ctx).Where("is_delete = ?", false).Find(&albums).Error
	if err != nil {
		return nil, err
	}
	return albums, nil
}

// GetPhotoAlbumByID 查询相册信息
func (dao *photoAlbumDao) GetPhotoAlbumByID(ctx context.Context, albumID int) (*model.PhotoAlbum, error) {
	var album model.PhotoAlbum
	err := dao.db.WithContext(ctx).First(&album, albumID).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

// LogicalDeletePhotoAlbum 逻辑删除相册
func (dao *photoAlbumDao) LogicalDeletePhotoAlbum(ctx context.Context, albumID int) error {
	return dao.db.WithContext(ctx).Model(&model.PhotoAlbum{}).Where("id = ?", albumID).Update("is_delete", true).Error
}

// DeletePhotoAlbum 删除相册
func (dao *photoAlbumDao) DeletePhotoAlbum(ctx context.Context, albumID int) error {
	return dao.db.WithContext(ctx).Delete(&model.PhotoAlbum{}, albumID).Error
}

// ListPublishedPhotoAlbums 查询相册列表（发布状态）
func (dao *photoAlbumDao) ListPublishedPhotoAlbums(ctx context.Context) ([]model.PhotoAlbum, error) {
	var albums []model.PhotoAlbum
	err := dao.db.WithContext(ctx).Where("status = ? AND is_delete = ?", "PUBLIC", false).Order("id DESC").Find(&albums).Error
	if err != nil {
		return nil, err
	}
	return albums, nil
}

// SavePhotoAlbum 保存相册
func (dao *photoAlbumDao) SavePhotoAlbum(ctx context.Context, album *model.PhotoAlbum) error {
	now := time.Now()
	album.CreateTime = now
	album.UpdateTime = nil // 创建时不设置更新时间
	return dao.db.WithContext(ctx).Create(album).Error
}

// UpdatePhotoAlbum 更新相册
func (dao *photoAlbumDao) UpdatePhotoAlbum(ctx context.Context, album *model.PhotoAlbum) error {
	now := time.Now()
	album.UpdateTime = &now
	return dao.db.WithContext(ctx).Model(&model.PhotoAlbum{}).Where("id = ?", album.ID).Omit("create_time").Updates(album).Error
}

// CountPhotosInAlbum 查询相册中的照片数量根据Id
func (dao *photoAlbumDao) CountPhotosInAlbumById(ctx context.Context, albumID int) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.Photo{}).Where("album_id = ? AND is_delete = ?", albumID, false).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdatePhotoAlbumByID 根据ID更新相册
func (dao *photoAlbumDao) UpdatePhotoAlbumByID(ctx context.Context, albumID int, updates map[string]interface{}) error {
	return dao.db.WithContext(ctx).Model(&model.PhotoAlbum{}).Where("id = ?", albumID).Updates(updates).Error
}

func (dao *photoAlbumDao) UpdateBatchByID(ctx context.Context, albums []model.PhotoAlbum) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		for _, album := range albums {
			if err := tx.Model(&album).Updates(map[string]interface{}{
				"is_delete":   album.IsDelete,
				"update_time": now,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (dao *photoAlbumDao) FindAlbumByID(ctx context.Context, albumID int, isDelete int, status int) (*model.PhotoAlbum, error) {
	var photoAlbum model.PhotoAlbum
	err := dao.db.WithContext(ctx).Where("id = ? AND is_delete = ? AND status = ?", albumID, isDelete, status).First(&photoAlbum).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &photoAlbum, nil
}
