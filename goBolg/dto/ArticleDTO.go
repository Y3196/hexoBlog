package dto

import "time"

// ArticleDTO 代表文章的 DTO
type ArticleDTO struct {
	ID                   int                     `json:"id"`
	UserID               int                     `json:"user_id"`
	ArticleCover         string                  `json:"articleCover"`
	ArticleTitle         string                  `json:"articleTitle"`
	ArticleContent       string                  `json:"articleContent"`
	LikeCount            int                     `json:"likeCount"`
	ViewsCount           int                     `json:"viewsCount"`
	Type                 int                     `json:"type"`
	OriginalUrl          string                  `json:"originalUrl"`
	CreateTime           time.Time               `json:"createTime"`
	UpdateTime           time.Time               `json:"updateTime"`
	CategoryId           int                     `json:"categoryId"`
	CategoryName         string                  `json:"categoryName"`
	TagDTOList           []TagDTO                `json:"tagDTOList" gorm:"many2many:tb_article_tag"`
	LastArticle          *ArticlePaginationDTO   `json:"lastArticle"`
	NextArticle          *ArticlePaginationDTO   `json:"nextArticle"`
	RecommendArticleList ArticleRecommendDTOList `json:"recommendArticleList"`
	NewestArticleList    ArticleRecommendDTOList `json:"newestArticleList"`
	UserLiked            bool                    `json:"userLiked"`
}
