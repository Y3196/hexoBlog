package vo

// RoleVO 代表角色信息
type RoleVO struct {
	ID             int    `json:"id"`                            // 角色id
	RoleName       string `json:"roleName" validate:"required"`  // 角色名
	RoleLabel      string `json:"roleLabel" validate:"required"` // 权限标签
	ResourceIDList []int  `json:"resourceIdList"`                // 资源id列表
	MenuIDList     []int  `json:"menuIdList"`                    // 菜单id列表
}
