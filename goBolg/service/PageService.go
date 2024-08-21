package service

import (
	"context"
	"goBolg/vo"
)

type PageService interface {
	SaveOrUpdatePage(ctx context.Context, pageVO vo.PageVO) error

	DeletePage(ctx context.Context, pageID uint) error

	ListPages(ctx context.Context) ([]vo.PageVO, error)
}
