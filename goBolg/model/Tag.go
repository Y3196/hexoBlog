package model

import (
	"time"
)

// Tag represents the tag entity
type Tag struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TagName    string    `gorm:"type:varchar(255);not null" json:"tagName"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

// TableName 设置表名
func (Tag) TableName() string {
	return "tb_tag"
}
