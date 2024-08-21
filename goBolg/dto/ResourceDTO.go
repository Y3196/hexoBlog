package dto

import "time"

// ResourceDTO 代表资源的 DTO
type ResourceDTO struct {
	ID            uint          `json:"id"`            // 权限id
	ResourceName  string        `json:"resourceName"`  // 资源名
	URL           string        `json:"url"`           // 权限路径
	RequestMethod string        `json:"requestMethod"` // 请求方式
	IsDisable     int           `json:"isDisable"`     // 是否禁用
	IsAnonymous   int           `json:"isAnonymous"`   // 是否匿名访问
	CreateTime    time.Time     `json:"createTime"`    // 创建时间
	Children      []ResourceDTO `json:"children"`      // 权限列表
}
