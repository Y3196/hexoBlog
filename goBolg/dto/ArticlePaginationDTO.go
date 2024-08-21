package dto

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ArticlePaginationDTO 代表文章分页的 DTO
type ArticlePaginationDTO struct {
	ID           int       `json:"id"`
	ArticleTitle string    `json:"articleTitle"`
	ArticleCover string    `json:"articleCover"`
	CreateTime   time.Time `json:"createTime"`
}

// 实现 Valuer 接口，将结构体转换为数据库值
func (a ArticlePaginationDTO) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// 实现 Scanner 接口，从数据库值读取结构体
func (a *ArticlePaginationDTO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, a)
}
