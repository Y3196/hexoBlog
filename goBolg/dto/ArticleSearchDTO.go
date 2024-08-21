package dto

type ArticleSearchDTO struct {
	ID             int    `json:"id"`             // 文章 id
	ArticleTitle   string `json:"articleTitle"`   // 文章标题
	ArticleContent string `json:"articleContent"` // 文章内容
	IsDelete       int    `json:"isDelete"`       // 是否删除
	Status         int    `json:"status"`         // 文章状态
}
