package Impl

import (
	"context"
	"github.com/jinzhu/copier"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/exception"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"gorm.io/gorm"
	"log"
)

type photoAlbumServiceImpl struct {
	photoAlbumDao dao.PhotoAlbumDao
	photoDao      dao.PhotoDao
}

// NewPhotoAlbumService 创建新的 PhotoAlbumService 实例
func NewPhotoAlbumService(photoAlbumDao dao.PhotoAlbumDao, photoDao dao.PhotoDao) service.PhotoAlbumService {
	return &photoAlbumServiceImpl{
		photoAlbumDao: photoAlbumDao,
		photoDao:      photoDao,
	}
}

func (service *photoAlbumServiceImpl) SaveOrUpdatePhotoAlbum(ctx context.Context, photoAlbumVO vo.PhotoAlbumVO) error {
	// 查询相册是否存在
	var photoAlbum model.PhotoAlbum
	if err := copier.Copy(&photoAlbum, &photoAlbumVO); err != nil {
		return err
	}

	if photoAlbumVO.ID != 0 {
		existingAlbum, err := service.photoAlbumDao.GetPhotoAlbumByID(ctx, photoAlbumVO.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existingAlbum != nil {
			// 更新
			if err := service.photoAlbumDao.UpdatePhotoAlbum(ctx, &photoAlbum); err != nil {
				return err
			}
		} else {
			// 插入
			if err := service.photoAlbumDao.SavePhotoAlbum(ctx, &photoAlbum); err != nil {
				return err
			}
		}
	} else {
		// 查询相册名是否存在
		existingAlbum, err := service.photoAlbumDao.FindAlbumByName(ctx, photoAlbumVO.AlbumName)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existingAlbum != nil {
			return exception.NewBizException("相册名已存在")
		}

		// 保存
		if err := service.photoAlbumDao.SavePhotoAlbum(ctx, &photoAlbum); err != nil {
			return err
		}
	}

	log.Println("Photo album saved or updated successfully")
	return nil
}

func (service *photoAlbumServiceImpl) ListPhotoAlbumBacks(ctx context.Context, condition vo.ConditionVO) (*vo.PageResult[dto.PhotoAlbumBackDTO], error) {
	// 查询相册数量
	count, err := service.photoAlbumDao.CountPhotoAlbums(ctx, condition)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return &vo.PageResult[dto.PhotoAlbumBackDTO]{}, nil
	}

	// 获取分页参数
	current := utils.GetLimitCurrent(ctx)
	size := utils.GetSize(ctx)

	// 查询相册信息
	photoAlbumBackList, err := service.photoAlbumDao.ListPhotoAlbumBacks(ctx, current, size, condition)
	if err != nil {
		return nil, err
	}

	return &vo.PageResult[dto.PhotoAlbumBackDTO]{
		RecordList: photoAlbumBackList,
		Count:      count,
	}, nil
}

// ListPhotoAlbumBackInfos 获取相册信息列表
func (service *photoAlbumServiceImpl) ListPhotoAlbumBackInfos(ctx context.Context) ([]dto.PhotoAlbumDTO, error) {
	photoAlbumList, err := service.photoAlbumDao.ListPhotoAlbums(ctx)
	if err != nil {
		return nil, err
	}

	var photoAlbumDTOs []dto.PhotoAlbumDTO
	err = copier.Copy(&photoAlbumDTOs, &photoAlbumList)
	if err != nil {
		return nil, err
	}

	return photoAlbumDTOs, nil
}

// GetPhotoAlbumBackByID 获取相册详细信息
func (service *photoAlbumServiceImpl) GetPhotoAlbumBackByID(ctx context.Context, albumID int) (*dto.PhotoAlbumBackDTO, error) {
	photoAlbum, err := service.photoAlbumDao.GetPhotoAlbumByID(ctx, albumID)
	if err != nil {
		return nil, err
	}

	photoCount, err := service.photoAlbumDao.CountPhotosInAlbumById(ctx, albumID)
	if err != nil {
		return nil, err
	}

	var albumDTO dto.PhotoAlbumBackDTO
	err = copier.Copy(&albumDTO, &photoAlbum)
	if err != nil {
		return nil, err
	}
	albumDTO.PhotoCount = int(photoCount)

	return &albumDTO, nil
}

func (service *photoAlbumServiceImpl) DeletePhotoAlbumByID(ctx context.Context, albumID int) error {
	count, err := service.photoDao.CountPhotos(ctx, albumID)
	if err != nil {
		return err
	}

	if count > 0 {
		updates := map[string]interface{}{
			"is_delete": true,
		}
		if err := service.photoAlbumDao.UpdatePhotoAlbumByID(ctx, albumID, updates); err != nil {
			return err
		}
		if err := service.photoDao.LogicalDeletePhotos(ctx, albumID); err != nil {
			return err
		}
		log.Println("Photo album and photos logically deleted successfully")
	} else {
		if err := service.photoAlbumDao.DeletePhotoAlbum(ctx, albumID); err != nil {
			return err
		}
		log.Println("Photo album physically deleted successfully")
	}

	return nil
}

func (service *photoAlbumServiceImpl) ListPhotoAlbums(ctx context.Context) ([]dto.PhotoAlbumDTO, error) {
	// 查询相册列表
	photoAlbumList, err := service.photoAlbumDao.ListPhotoAlbums(ctx)
	if err != nil {
		return nil, err
	}
	// 复制列表
	var photoAlbumDTOList []dto.PhotoAlbumDTO
	err = copier.Copy(&photoAlbumDTOList, &photoAlbumList)
	if err != nil {
		return nil, err
	}
	return photoAlbumDTOList, nil
}
