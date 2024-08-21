package dto

import (
	"time"
)

// ArticlePreviewDTO 代表预览文章 DTO
type ArticlePreviewDTO struct {
	ID           int       `json:"id"`           // 文章 id
	ArticleCover string    `json:"articleCover"` // 文章缩略图
	ArticleTitle string    `json:"articleTitle"` // 标题
	CreateTime   time.Time `json:"createTime"`   // 发表时间
	CategoryId   int       `json:"categoryId"`   // 文章分类 id
	CategoryName string    `json:"categoryName"` // 文章分类名
	TagDTOList   []TagDTO  `json:"tagDTOList"`   // 文章标签
}
