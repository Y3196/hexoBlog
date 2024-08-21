package vo

import (
	"time"
)

// HexoArticleVO 代表 Hexo 文章，继承自 ArticleVO
type HexoArticleVO struct {
	ArticleVO
	CreateTime time.Time `json:"createTime"` // 创建时间
}
