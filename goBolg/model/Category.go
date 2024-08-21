package model

import (
	"time"
)

// Category represents a category entity
type Category struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string    `gorm:"column:category_name"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (Category) TableName() string {
	return "tb_category"
}
