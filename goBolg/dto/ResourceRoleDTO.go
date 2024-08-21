package dto

// ResourceRoleDTO 代表资源角色的 DTO
type ResourceRoleDTO struct {
	ID            int      `json:"id"`            // 资源id
	URL           string   `json:"url"`           // 路径
	RequestMethod string   `json:"requestMethod"` // 请求方式
	RoleList      []string `json:"roleList"`      // 角色名
}
