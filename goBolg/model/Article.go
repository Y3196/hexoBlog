package model

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID             int        `gorm:"primaryKey"`
	UserID         int        `gorm:"column:user_id"`
	CategoryID     int        `gorm:"column:category_id"`
	ArticleCover   string     `gorm:"column:article_cover"`
	ArticleTitle   string     `gorm:"column:article_title"`
	ArticleContent string     `gorm:"column:article_content"`
	Type           *int       `gorm:"column:type"`
	OriginalURL    *string    `gorm:"column:original_url"`
	IsTop          *int       `gorm:"column:is_top"`
	IsDelete       *int       `gorm:"column:is_delete;default:0"`
	Status         *int       `gorm:"column:status"`
	CreateTime     time.Time  `gorm:"column:create_time;autoCreateTime"`
	UpdateTime     *time.Time `gorm:"column:update_time"`
}

func (Article) TableName() string {
	return "tb_article"
}

// BeforeCreate 在创建记录之前设置创建时间
func (a *Article) BeforeCreate(tx *gorm.DB) (err error) {
	a.CreateTime = time.Now()
	a.UpdateTime = nil // 确保在创建时不设置更新时间
	return
}

// BeforeUpdate 在更新记录之前设置更新时间
func (a *Article) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	a.UpdateTime = &now
	return
}
