package dto

// CommentCountDTO 代表评论数量 DTO
type CommentCountDTO struct {
	ID           int `json:"id"`           // 评论 ID
	CommentCount int `json:"commentCount"` // 评论数量
}
