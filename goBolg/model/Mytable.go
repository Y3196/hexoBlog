package model

import (
	"fmt"
	"time"
)

// Mytable 代表表的实体
type Mytable struct {
	// 主键
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 信息
	Info string `json:"info" gorm:"column:info"`

	// 创建时间（如果需要）
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime;column:create_time"`

	// 修改时间（如果需要）
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime;column:update_time"`
}

// String 返回结构体的字符串表示
func (mt Mytable) String() string {
	return "Mytable{" +
		"id=" + fmt.Sprint(mt.ID) +
		", info='" + mt.Info + "'" +
		"}"
}
