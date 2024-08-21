package dto

import "time"

type ArchiveDTO struct {
	ID           uint      `json:"id"`
	ArticleTitle string    `json:"article_title"`
	CreateTime   time.Time `json:"create_time"`
}
