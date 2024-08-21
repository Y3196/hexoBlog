package dto

import (
	"time"
)

// ArticleBackDTO 代表后台文章的 DTO
type ArticleBackDTO struct {
	ID           int       `json:"id"`
	ArticleCover string    `json:"articleCover"`
	ArticleTitle string    `json:"articleTitle"`
	CreateTime   time.Time `json:"createTime"`
	LikeCount    int       `json:"likeCount"`
	ViewsCount   int       `json:"viewsCount"`
	CategoryName string    `json:"categoryName"`
	TagDTOList   []TagDTO  `json:"tagDTOList"`
	Type         int       `json:"type"`
	IsTop        int       `json:"isTop"`
	IsDelete     int       `json:"isDelete"`
	Status       int       `json:"status"`
}
