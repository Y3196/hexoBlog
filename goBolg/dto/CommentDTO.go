package dto

import "time"

// CommentDTO 代表评论 DTO
type CommentDTO struct {
	ID             int        `json:"id"`                                                    // 评论 ID
	UserID         int        `json:"userId"`                                                // 用户 ID
	Nickname       string     `json:"nickname"`                                              // 用户昵称
	Avatar         string     `json:"avatar"`                                                // 用户头像
	WebSite        string     `json:"webSite"`                                               // 个人网站
	CommentContent string     `json:"commentContent"`                                        // 评论内容
	LikeCount      int        `json:"likeCount"`                                             // 点赞数
	CreateTime     time.Time  `json:"createTime"`                                            // 评论时间
	ReplyCount     int        `json:"replyCount"`                                            // 回复量
	ReplyDTOList   []ReplyDTO `json:"replyDTOList" gorm:"foreignKey:ParentID;references:ID"` // 回复列表
}
