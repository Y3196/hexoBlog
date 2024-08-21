package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
	"net/http"
)

// BlogInfoService defines the rabbitService for fetching blog information.
type BlogInfoService interface {
	GetBlogHomeInfo(ctx context.Context) (dto.BlogHomeInfoDTO, error)

	// GetBlogBackInfo retrieves the backend information of the blog.
	GetBlogBackInfo(ctx context.Context) (dto.BlogBackInfoDTO, error)

	// UpdateWebsiteConfig saves or updates the website configuration.
	UpdateWebsiteConfig(ctx context.Context, websiteConfigVO vo.WebsiteConfigVO) error

	// GetWebsiteConfig retrieves the current website configuration.
	GetWebsiteConfig(ctx context.Context) (vo.WebsiteConfigVO, error)

	// GetAbout retrieves the "About Me" content.
	GetAbout(ctx context.Context) (string, error)

	// UpdateAbout updates the "About Me" content.
	UpdateAbout(ctx context.Context, blogInfoVO vo.BlogInfoVO) error

	// Report uploads visitor information.
	Report(ctx context.Context, request *http.Request)
}
