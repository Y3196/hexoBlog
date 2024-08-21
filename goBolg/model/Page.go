package model

import (
	"time"
)

// Page 页面实体
type Page struct {
	ID         int        `gorm:"primaryKey;autoIncrement"`
	PageName   string     `gorm:"column:page_name"`
	PageLabel  string     `gorm:"column:page_label"`
	PageCover  string     `gorm:"column:page_cover"`
	CreateTime time.Time  `gorm:"column:create_time;autoCreateTime"`
	UpdateTime *time.Time `gorm:"column:update_time"`
}

func (Page) TableName() string {
	return "tb_page"
}
