package model

import (
	"time"
)

// Message 表示留言数据库表的结构。
type Message struct {
	ID             int       `gorm:"primaryKey;autoIncrement"` // 主键ID
	IPAddress      string    `gorm:"size:100"`                 // 用户IP
	IPSource       string    `gorm:"size:100"`                 // 用户地址
	Nickname       string    `gorm:"size:100"`                 // 昵称
	Avatar         string    `gorm:"size:255"`                 // 头像
	MessageContent string    `gorm:"type:text"`                // 留言内容
	Time           int       `gorm:""`                         // 弹幕速度
	IsReview       int       `gorm:"default:1"`                // 是否审核
	CreateTime     time.Time `gorm:"autoCreateTime"`           // 创建时间
	UpdateTime     time.Time `gorm:"column:update_time"`
}

func (Message) TableName() string {
	return "tb_message"
}
