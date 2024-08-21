package dto

import "time"

// ReplyDTO 代表回复的 DTO
type ReplyDTO struct {
	ID             int       `json:"id"`             // 评论id
	ParentID       int       `json:"parentId"`       // 父评论id
	UserID         int       `json:"userId"`         // 用户id
	Nickname       string    `json:"nickname"`       // 用户昵称
	Avatar         string    `json:"avatar"`         // 用户头像
	WebSite        string    `json:"webSite"`        // 个人网站
	ReplyUserID    int       `json:"replyUserId"`    // 被回复用户id
	ReplyNickname  string    `json:"replyNickname"`  // 被回复用户昵称
	ReplyWebSite   string    `json:"replyWebSite"`   // 被回复个人网站
	CommentContent string    `json:"commentContent"` // 评论内容
	LikeCount      int       `json:"likeCount"`      // 点赞数
	CreateTime     time.Time `json:"createTime"`     // 评论时间
}
