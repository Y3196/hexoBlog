package dto

import "time"

// CategoryBackDTO 代表后台分类 DTO
type CategoryBackDTO struct {
	ID           int       `json:"id"`           // 分类 id
	CategoryName string    `json:"categoryName"` // 分类名称
	ArticleCount int       `json:"articleCount"` // 文章数量
	CreateTime   time.Time `json:"createTime"`   // 创建时间
}
