package dto

import "time"

type ArticleHomeDTO struct {
	ID           int    `json:"id"`
	ArticleCover string `json:"articleCover"`

	Tags           []TagDTO  `json:"tags"`
	ArticleTitle   string    `json:"articleTitle"`
	ArticleContent string    `json:"articleContent"`
	CreateTime     time.Time `json:"createTime"`
	IsTop          int       `json:"isTop"`
	Type           int       `json:"type"`
	CategoryID     int       `json:"categoryId"`
	CategoryName   string    `json:"categoryName"`
	TagDTOList     []TagDTO  `json:"tagDTOList" gorm:"-"`
}
