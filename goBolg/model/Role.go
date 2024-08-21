package model

import (
	"fmt"
	"time"
)

// Role 代表角色实体
type Role struct {
	// 角色id
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 角色名
	RoleName string `json:"roleName" gorm:"column:role_name"`

	// 角色标签
	RoleLabel string `json:"roleLabel" gorm:"column:role_label"`

	// 是否禁用
	IsDisable int `json:"isDisable" gorm:"column:is_disable"`

	// 创建时间
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime;column:create_time"`

	// 修改时间
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime;column:update_time"`
}

// String 返回结构体的字符串表示
func (r Role) String() string {
	return fmt.Sprintf("Role{ID=%d, RoleName='%s', RoleLabel='%s', IsDisable=%d, CreateTime=%s, UpdateTime=%s}",
		r.ID, r.RoleName, r.RoleLabel, r.IsDisable, r.CreateTime.Format(time.RFC3339), r.UpdateTime.Format(time.RFC3339))
}

func (Role) TableName() string {
	return "tb_role"
}
