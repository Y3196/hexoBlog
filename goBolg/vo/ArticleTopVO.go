package vo

// ArticleTopVO 代表文章置顶信息
type ArticleTopVO struct {
	ID    *int `json:"id" validate:"required"`     // id，必填
	IsTop *int `json:"is_top" validate:"required"` // 置顶状态，必填
}
