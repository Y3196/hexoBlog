package dto

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ArticleRecommendDTO 代表推荐文章 DTO
type ArticleRecommendDTO struct {
	ID           int       `json:"id"`           // 文章 id
	ArticleCover string    `json:"articleCover"` // 文章缩略图
	ArticleTitle string    `json:"articleTitle"` // 标题
	CreateTime   time.Time `json:"createTime"`   // 创建时间
}

// ArticleRecommendDTOList 代表推荐文章的 DTO 列表
type ArticleRecommendDTOList []ArticleRecommendDTO

// 实现 Valuer 接口，将结构体转换为数据库值
func (a ArticleRecommendDTOList) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// 实现 Scanner 接口，从数据库值读取结构体
func (a *ArticleRecommendDTOList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, a)
}
