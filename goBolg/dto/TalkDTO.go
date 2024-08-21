package dto

import "time"

// TalkDTO 代表说说 DTO
type TalkDTO struct {
	ID           int       `json:"id"` // 说说id
	UserID       int       `json:"user_id"`
	Nickname     string    `json:"nickname"`     // 昵称
	Avatar       string    `json:"avatar"`       // 头像
	Content      string    `json:"content"`      // 说说内容
	Images       string    `json:"images"`       // 图片
	ImgList      []string  `json:"imgList"`      // 图片列表
	IsTop        int       `json:"isTop"`        // 是否置顶
	LikeCount    int       `json:"likeCount"`    // 点赞量
	CommentCount int       `json:"commentCount"` // 评论量
	CreateTime   time.Time `json:"createTime"`   // 创建时间
}
