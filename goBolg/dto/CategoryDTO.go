package dto

// CategoryDTO 代表了一个文章分类的数据。
type CategoryDTO struct {
	ID           uint   `json:"id"`           // ID 是分类的唯一标识符。
	CategoryName string `json:"categoryName"` // CategoryName 是分类的名称。
	ArticleCount int    `json:"articleCount"` // ArticleCount 是属于这个分类的文章数量。
}
