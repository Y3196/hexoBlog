package model

import (
	"time"

	"gorm.io/gorm"
)

// ChatRecord 聊天记录
type ChatRecord struct {
	gorm.Model // 自动包含 ID、CreatedAt、UpdatedAt 和 DeletedAt

	// 用户ID
	UserID int `gorm:"column:user_id" json:"user_id"`

	// 用户昵称
	Nickname string `gorm:"column:nickname" json:"nickname"`

	// 用户头像
	Avatar string `gorm:"column:avatar" json:"avatar"`

	// 聊天内容
	Content string `gorm:"column:content" json:"content"`

	// 类型
	Type int `gorm:"column:type" json:"type"`

	// 用户登录IP
	IPAddress string `gorm:"column:ip_address" json:"ip_address"`

	// IP来源
	IPSource string `gorm:"column:ip_source" json:"ip_source"`

	// 创建时间
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`

	// 修改时间
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 设置表名
func (ChatRecord) TableName() string {
	return "tb_chat_record"
}
