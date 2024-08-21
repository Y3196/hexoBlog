package model

import (
	"fmt"
	"time"
)

// Talk 代表说说实体
type Talk struct {
	// 说说id
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 用户id
	UserID int `json:"userId" gorm:"column:user_id"`

	// 说说内容
	Content string `json:"content" gorm:"column:content"`

	// 图片
	Images string `json:"images" gorm:"column:images"`

	// 是否置顶
	IsTop int `json:"isTop" gorm:"column:is_top"`

	// 说说状态 1.公开 2.私密
	Status int `json:"status" gorm:"column:status"`

	// 创建时间
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime;column:create_time"`

	// 修改时间
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime;column:update_time"`
}

// String 返回结构体的字符串表示
func (t Talk) String() string {
	return fmt.Sprintf("Talk{ID=%d, UserID=%d, Content='%s', Images='%s', IsTop=%d, Status=%d, CreateTime=%s, UpdateTime=%s}",
		t.ID, t.UserID, t.Content, t.Images, t.IsTop, t.Status, t.CreateTime.Format(time.RFC3339), t.UpdateTime.Format(time.RFC3339))
}

// TableName 设置表名
func (Talk) TableName() string {
	return "tb_talk"
}
