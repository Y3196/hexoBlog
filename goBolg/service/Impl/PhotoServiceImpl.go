package Impl

import (
	"context"
	"github.com/bwmarrin/snowflake"
	constants "goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/enums"
	"goBolg/exception"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"log"
	"time"

	"github.com/jinzhu/copier"
)

type photoServiceImpl struct {
	photoDao      dao.PhotoDao
	node          *snowflake.Node
	photoAlbumDao dao.PhotoAlbumDao
}

// NewPhotoService 创建新的 PhotoService 实例
func NewPhotoServiceImpl(photoDao dao.PhotoDao, photoAlbumDao dao.PhotoAlbumDao) (service.PhotoService, error) {
	node, err := snowflake.NewNode(1) // 这里的参数可以是任意你喜欢的数字，表示节点ID
	if err != nil {
		log.Fatalf("Failed to initialize snowflake node: %v", err)
		return nil, err
	}
	return &photoServiceImpl{photoDao: photoDao, photoAlbumDao: photoAlbumDao, node: node}, nil
}

func (service *photoServiceImpl) ListPhotos(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.PhotoBackDTO], error) {
	current := utils.GetCurrent(ctx)
	size := utils.GetSize(ctx)

	photos, total, err := service.photoDao.ListPhotos(ctx, current, size, condition)
	if err != nil {
		return vo.PageResult[dto.PhotoBackDTO]{}, err
	}

	var photoBackDTOList []dto.PhotoBackDTO
	err = copier.Copy(&photoBackDTOList, &photos)
	if err != nil {
		return vo.PageResult[dto.PhotoBackDTO]{}, err
	}

	return vo.NewPageResult(photoBackDTOList, total), nil
}

// UpdatePhoto 更新照片信息
func (s *photoServiceImpl) UpdatePhoto(ctx context.Context, photoInfoVO vo.PhotoInfoVO) error {
	photo := &model.Photo{}
	utils.BeanCopy(photoInfoVO, photo)
	return s.photoDao.UpdateByID(ctx, photo)
}

func (s *photoServiceImpl) SavePhotos(ctx context.Context, photoVO vo.PhotoVO) error {
	var photoList []model.Photo
	for _, item := range photoVO.PhotoURLList {
		photoList = append(photoList, model.Photo{
			AlbumID:    photoVO.AlbumID,
			PhotoName:  s.node.Generate().String(), // 使用雪花算法生成ID
			PhotoSrc:   item,
			CreateTime: time.Now(),
		})
	}

	return s.photoDao.SaveBatch(ctx, photoList)
}

func (s *photoServiceImpl) UpdatePhotosAlbum(ctx context.Context, photoVO vo.PhotoVO) error {
	var photos []model.Photo
	for _, id := range photoVO.PhotoIDList {
		photos = append(photos, model.Photo{
			ID:      id,
			AlbumID: photoVO.AlbumID,
		})
	}
	err := s.photoDao.UpdateBatch(ctx, photos)
	if err != nil {
		log.Printf("Error updating photos album: %v", err)
		return err
	}
	return nil
}

func (s *photoServiceImpl) UpdatePhotoDelete(ctx context.Context, deleteVO vo.DeleteVO) error {
	var photos []model.Photo
	for _, id := range deleteVO.IDList {
		photos = append(photos, model.Photo{
			ID:       id,
			IsDelete: *deleteVO.IsDelete,
		})
	}

	err := s.photoDao.UpdateBatchByID(ctx, photos)
	if err != nil {
		log.Printf("Error updating photos delete status: %v", err)
		return err
	}

	if *deleteVO.IsDelete == constants.False {
		photoList, err := s.photoDao.SelectListByIDs(ctx, deleteVO.IDList)
		if err != nil {
			log.Printf("Error selecting photo list: %v", err)
			return err
		}

		var albumIDs []int
		for _, photo := range photoList {
			albumIDs = append(albumIDs, photo.AlbumID)
		}

		photoAlbums := make(map[int]struct{})
		for _, albumID := range albumIDs {
			if _, exists := photoAlbums[albumID]; !exists {
				photoAlbums[albumID] = struct{}{}
			}
		}

		var photoAlbumList []model.PhotoAlbum
		for albumID := range photoAlbums {
			photoAlbumList = append(photoAlbumList, model.PhotoAlbum{
				ID:       albumID,
				IsDelete: constants.False,
			})
		}

		err = s.photoAlbumDao.UpdateBatchByID(ctx, photoAlbumList)
		if err != nil {
			log.Printf("Error updating photo albums delete status: %v", err)
			return err
		}
	}
	return nil
}

func (s *photoServiceImpl) DeletePhotos(ctx context.Context, photoIdList []int) error {
	err := s.photoDao.DeleteBatchByIDs(ctx, photoIdList)
	if err != nil {
		log.Printf("Error deleting photos: %v", err)
		return err
	}
	return nil
}

func (s *photoServiceImpl) ListPhotosByAlbumID(ctx context.Context, albumID int, current int, size int) (*dto.PhotoDTO, error) {
	// 查询相册信息
	photoAlbum, err := s.photoAlbumDao.FindAlbumByID(ctx, albumID, constants.False, enums.PUBLIC_STATUS.Status)
	if err != nil {
		log.Printf("Error finding album: %v", err)
		return nil, err
	}
	if photoAlbum == nil {
		return nil, exception.NewBizException("相册不存在")
	}

	// 查询照片列表
	photos, err := s.photoDao.SelectPhotosByAlbumID(ctx, albumID, current, size)
	if err != nil {
		log.Printf("Error selecting photos: %v", err)
		return nil, err
	}

	photoSrcList := make([]string, len(photos))
	for i, photo := range photos {
		photoSrcList[i] = photo.PhotoSrc
	}

	return &dto.PhotoDTO{
		PhotoAlbumCover: photoAlbum.AlbumCover,
		PhotoAlbumName:  photoAlbum.AlbumName,
		PhotoList:       photoSrcList,
	}, nil
}
