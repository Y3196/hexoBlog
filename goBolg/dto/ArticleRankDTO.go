package dto

// ArticleRankDTO represents the ranking data for articles.
type ArticleRankDTO struct {
	ArticleTitle string `json:"articleTitle"` // ArticleTitle is the title of the article.
	ViewsCount   int    `json:"viewsCount"`   // ViewsCount is the number of views for the article.
}
