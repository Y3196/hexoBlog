package dto

import "time"

// CommentBackDTO 代表后台评论 DTO
type CommentBackDTO struct {
	ID             int       `json:"id"`             // 评论 ID
	Avatar         string    `json:"avatar"`         // 用户头像
	Nickname       string    `json:"nickname"`       // 用户昵称
	ReplyNickname  string    `json:"replyNickname"`  // 被回复用户昵称
	ArticleTitle   string    `json:"articleTitle"`   // 文章标题
	CommentContent string    `json:"commentContent"` // 评论内容
	Type           int       `json:"type"`           // 评论类型
	IsReview       int       `json:"isReview"`       // 是否审核
	CreateTime     time.Time `json:"createTime"`     // 发表时间
}
