package model

import (
	"gorm.io/gorm"
	"time"
)

// Menu 代表菜单实体
type Menu struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 菜单名
	Name string `json:"name" gorm:"column:name"`

	// 路径
	Path string `json:"path" gorm:"column:path"`

	// 组件
	Component string `json:"component" gorm:"column:component"`

	// icon
	Icon string `json:"icon" gorm:"column:icon"`

	// 排序
	OrderNum int `json:"orderNum" gorm:"column:order_num"`

	// 父id
	ParentID *uint `json:"parentId" gorm:"column:parent_id"`

	// 是否隐藏
	IsHidden *int `json:"isHidden" gorm:"column:is_hidden"`

	// 创建时间
	CreateTime time.Time `gorm:"column:create_time"`

	// 修改时间
	UpdateTime *time.Time `gorm:"column:update_time"`
}

func (Menu) TableName() string {
	return "tb_menu"
}

// BeforeCreate 在创建记录之前设置创建时间
func (m *Menu) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return
}

// BeforeUpdate 在更新记录之前设置更新时间
func (m *Menu) BeforeUpdate(tx *gorm.DB) (err error) {
	t := time.Now()
	m.UpdateTime = &t
	return
}
