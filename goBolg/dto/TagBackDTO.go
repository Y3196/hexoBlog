package dto

import "time"

// TagBackDTO 代表后台标签
type TagBackDTO struct {
	ID           int       `json:"id"`           // 标签id
	TagName      string    `json:"tagName"`      // 标签名
	ArticleCount int       `json:"articleCount"` // 文章量
	CreateTime   time.Time `json:"createTime"`   // 创建时间
}
