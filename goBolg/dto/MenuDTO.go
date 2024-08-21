package dto

import (
	"time"
)

// MenuDTO 代表菜单 DTO
type MenuDTO struct {
	ID         int       `json:"id"`         // id
	Name       string    `json:"name"`       // 菜单名
	Path       string    `json:"path"`       // 路径
	Component  string    `json:"component"`  // 组件
	Icon       string    `json:"icon"`       // icon
	CreateTime time.Time `json:"createTime"` // 创建时间
	OrderNum   int       `json:"orderNum"`   // 排序
	IsDisable  int       `json:"isDisable"`  // 是否禁用
	IsHidden   int       `json:"isHidden"`   // 是否隐藏
	Children   []MenuDTO `json:"children"`   // 子菜单列表
}
