package dto

// UserRoleDTO 代表用户角色选项
type UserRoleDTO struct {
	ID       int    `json:"id"`       // 角色id
	RoleName string `json:"roleName"` // 角色名
}
