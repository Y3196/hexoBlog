package model

import (
	"time"

	"gorm.io/gorm"
)

// FriendLink 友链列表
type FriendLink struct {
	ID uint `gorm:"primaryKey"`

	// 链接名
	LinkName string `gorm:"column:link_name" json:"link_name"`

	// 链接头像
	LinkAvatar string `gorm:"column:link_avatar" json:"link_avatar"`

	// 链接地址
	LinkAddress string `gorm:"column:link_address" json:"link_address"`

	// 介绍
	LinkIntro string `gorm:"column:link_intro" json:"link_intro"`

	// 创建时间
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`

	// 修改时间
	UpdateTime *time.Time `json:"update_time" gorm:"column:update_time"`
}

// TableName 设置表名
func (FriendLink) TableName() string {
	return "tb_friend_link"
}

// BeforeCreate 在创建记录之前设置创建时间
func (f *FriendLink) BeforeCreate(tx *gorm.DB) (err error) {
	f.CreateTime = time.Now()
	f.UpdateTime = nil // 确保创建时不设置更新时间
	return
}

// BeforeUpdate 在更新记录之前设置更新时间
func (f *FriendLink) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	f.UpdateTime = &now
	return
}
