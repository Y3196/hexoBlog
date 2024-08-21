package model

import "fmt"

// RoleResource 代表角色资源实体
type RoleResource struct {
	// 主键id
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 角色id
	RoleID int `json:"roleId" gorm:"column:role_id"`

	// 资源id
	ResourceID int `json:"resourceId" gorm:"column:resource_id"`
}

// String 返回结构体的字符串表示
func (rr RoleResource) String() string {
	return fmt.Sprintf("RoleResource{ID=%d, RoleID=%d, ResourceID=%d}", rr.ID, rr.RoleID, rr.ResourceID)
}

// TableName 设置表名
func (RoleResource) TableName() string {
	return "tb_role_resource"
}
