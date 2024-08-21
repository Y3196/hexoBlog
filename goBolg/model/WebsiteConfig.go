package model

import (
	"time"
)

// WebsiteConfig represents the website_config table structure
type WebsiteConfig struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"` // 使用uint作为ID类型
	Config     string    `gorm:"type:varchar(255)"`        // 配置信息字段
	CreateTime time.Time `gorm:"autoCreateTime"`           // 创建时间，自动填充
	UpdateTime time.Time `gorm:"autoUpdateTime"`           // 修改时间，自动更新
}
