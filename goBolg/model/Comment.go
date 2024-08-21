package model

import (
	"time"
)

// Comment 评论
type Comment struct {
	ID int `gorm:"primaryKey"`

	// 评论用户ID
	UserID int `gorm:"column:user_id" json:"user_id"`

	// 回复用户ID
	ReplyUserID *int `gorm:"column:reply_user_id" json:"reply_user_id"`

	// 评论主题ID
	TopicID int `gorm:"column:topic_id" json:"topic_id"`

	// 评论内容
	CommentContent string `gorm:"column:comment_content" json:"comment_content"`

	// 父评论ID
	ParentID *int `gorm:"column:parent_id" json:"parent_id"`

	// 评论类型 1.文章 2.友链 3.说说
	Type int `gorm:"column:type" json:"type"`

	// 是否删除
	IsDelete int `gorm:"column:is_delete" json:"is_delete"`

	// 是否审核
	IsReview int `gorm:"column:is_review" json:"is_review"`

	// 创建时间
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`

	// 修改时间
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 设置表名
func (Comment) TableName() string {
	return "tb_comment"
}
