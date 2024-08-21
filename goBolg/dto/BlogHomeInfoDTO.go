package dto

import (
	"goBolg/vo"
)

type BlogHomeInfoDTO struct {
	ArticleCount  int                `json:"articleCount" description:"文章数量"`
	CategoryCount int                `json:"categoryCount" description:"分类数量"`
	TagCount      int                `json:"tagCount" description:"标签数量"`
	ViewsCount    string             `json:"viewsCount" description:"访问量"`
	WebsiteConfig vo.WebsiteConfigVO `json:"websiteConfig" description:"网站配置"`
	PageList      []vo.PageVO        `json:"pageList" description:"页面列表"`
}

// NewBlogHomeInfoDTO 创建一个新的 BlogHomeInfoDTO 实例
func NewBlogHomeInfoDTO(
	articleCount int,
	categoryCount int,
	tagCount int,
	viewsCount string,
	websiteConfig vo.WebsiteConfigVO,
	pageList []vo.PageVO,
) *BlogHomeInfoDTO {
	return &BlogHomeInfoDTO{
		ArticleCount:  articleCount,
		CategoryCount: categoryCount,
		TagCount:      tagCount,
		ViewsCount:    viewsCount,
		WebsiteConfig: websiteConfig,
		PageList:      pageList,
	}
}
