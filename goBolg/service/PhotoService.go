package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

type PhotoService interface {
	ListPhotos(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.PhotoBackDTO], error)

	UpdatePhoto(ctx context.Context, photoInfoVO vo.PhotoInfoVO) error

	SavePhotos(ctx context.Context, photoVO vo.PhotoVO) error

	UpdatePhotosAlbum(ctx context.Context, photoVO vo.PhotoVO) error

	UpdatePhotoDelete(ctx context.Context, deleteVO vo.DeleteVO) error

	DeletePhotos(ctx context.Context, photoIdList []int) error

	ListPhotosByAlbumID(ctx context.Context, albumID int, current int, size int) (*dto.PhotoDTO, error)
}
