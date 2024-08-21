package model

import "fmt"

// RoleMenu 代表角色菜单实体
type RoleMenu struct {
	// 主键id
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 角色id
	RoleID int `json:"roleId" gorm:"column:role_id"`

	// 菜单id
	MenuID int `json:"menuId" gorm:"column:menu_id"`
}

// String 返回结构体的字符串表示
func (rm RoleMenu) String() string {
	return fmt.Sprintf("RoleMenu{ID=%d, RoleID=%d, MenuID=%d}", rm.ID, rm.RoleID, rm.MenuID)
}
func (RoleMenu) TableName() string {
	return "tb_role_menu"
}
