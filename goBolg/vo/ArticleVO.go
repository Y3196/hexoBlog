package vo

import (
	"github.com/go-playground/validator/v10"
)

// ArticleVO 代表文章信息
type ArticleVO struct {
	ID             *int     `json:"id,omitempty"`                       // 文章id
	ArticleTitle   string   `json:"articleTitle" validate:"required"`   // 文章标题，必填
	ArticleContent string   `json:"articleContent" validate:"required"` // 文章内容，必填
	ArticleCover   *string  `json:"articleCover,omitempty"`             // 文章封面
	CategoryName   *string  `json:"categoryName,omitempty"`             // 文章分类
	TagNameList    []string `json:"tagNameList,omitempty"`              // 文章标签
	Type           *int     `json:"type,omitempty"`                     // 文章类型
	OriginalURL    *string  `json:"originalUrl,omitempty"`              // 原文链接
	IsTop          *int     `json:"isTop,omitempty"`                    // 是否置顶
	Status         *int     `json:"status,omitempty"`                   // 文章状态
}

// ValidateArticleVO 验证 ArticleVO 实例
func ValidateArticleVO(article ArticleVO) error {
	validate := validator.New()
	return validate.Struct(article)
}
