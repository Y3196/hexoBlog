package dto

// ReplyCountDTO 代表回复数量的 DTO
type ReplyCountDTO struct {
	CommentID  int `json:"commentId"`  // 评论id
	ReplyCount int `json:"replyCount"` // 回复数量
}
