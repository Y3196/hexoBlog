package model

import "fmt"

// UserRole 代表用户角色实体
type UserRole struct {
	// 主键id
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 用户id
	UserID int `json:"userId" gorm:"column:user_id"`

	// 角色id
	RoleID int `json:"roleId" gorm:"column:role_id"`
}

// String 返回结构体的字符串表示
func (ur UserRole) String() string {
	return fmt.Sprintf("UserRole{ID=%d, UserID=%d, RoleID=%d}", ur.ID, ur.UserID, ur.RoleID)
}

func (UserRole) TableName() string {
	return "tb_user_role"
}
