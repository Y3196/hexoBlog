package dto

// ArticlePreviewListDTO 代表文章预览列表 DTO
type ArticlePreviewListDTO struct {
	ArticlePreviewDTOList []ArticlePreviewDTO `json:"articlePreviewDTOList"` // 文章列表
	Name                  string              `json:"name"`                  // 条件名
}
