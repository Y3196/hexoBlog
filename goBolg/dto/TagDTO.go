package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// TagDTO 标签数据传输对象
type TagDTO struct {
	// ID 标签的唯一标识符
	ID int `json:"id"`

	// TagName 标签名
	TagName string `json:"tagName"`
}

// Tags 定义一个新类型来实现 Valuer 和 Scanner 接口
type Tags []TagDTO

// Value 实现 Valuer 接口
func (t Tags) Value() (driver.Value, error) {
	// 将 Tags 转换为 JSON 字符串
	return json.Marshal(t)
}

// Scan 实现 Scanner 接口
func (t *Tags) Scan(value interface{}) error {
	// 将数据库中的 JSON 字符串转换为 Tags
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, t)
}
