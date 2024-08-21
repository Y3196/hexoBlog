package model

import (
	"time"
)

// UniqueView represents the unique visits on the website.
type UniqueView struct {
	ID         int       `gorm:"primaryKey;autoIncrement"` // ID is the primary key and auto increments
	ViewsCount int       `gorm:"column:views_count"`       // Maps to views_count in the database
	CreateTime time.Time `gorm:"autoCreateTime"`           // Automatically handle creation time
	UpdateTime time.Time `gorm:"autoUpdateTime"`           // Automatically handle update time
}

// TableName overrides the table name used by User to `tb_unique_view`, which is the actual table name in the database.
func (UniqueView) TableName() string {
	return "tb_unique_view"
}
