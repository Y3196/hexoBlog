package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

// PhotoAlbumService 相册服务接口
type PhotoAlbumService interface {
	SaveOrUpdatePhotoAlbum(ctx context.Context, photoAlbumVO vo.PhotoAlbumVO) error

	ListPhotoAlbumBacks(ctx context.Context, condition vo.ConditionVO) (*vo.PageResult[dto.PhotoAlbumBackDTO], error)

	ListPhotoAlbumBackInfos(ctx context.Context) ([]dto.PhotoAlbumDTO, error)

	GetPhotoAlbumBackByID(ctx context.Context, albumID int) (*dto.PhotoAlbumBackDTO, error)

	DeletePhotoAlbumByID(ctx context.Context, albumID int) error

	ListPhotoAlbums(ctx context.Context) ([]dto.PhotoAlbumDTO, error)
}
