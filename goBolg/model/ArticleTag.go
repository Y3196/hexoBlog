package model

import (
	"gorm.io/gorm"
)

// ArticleTag 文章标签
type ArticleTag struct {
	gorm.Model // 自动包含 ID、CreatedAt、UpdatedAt 和 DeletedAt

	// 文章ID
	ArticleID int `gorm:"column:article_id" json:"article_id"`

	// 标签ID
	TagID int `gorm:"column:tag_id" json:"tag_id"`
}

// TableName 设置表名
func (ArticleTag) TableName() string {
	return "tb_article_tag"
}
