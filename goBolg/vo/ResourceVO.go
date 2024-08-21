package vo

// ResourceVO 代表资源信息
type ResourceVO struct {
	ID            int    `json:"id" validate:"required"`            // 资源id
	ResourceName  string `json:"resourceName" validate:"required"`  // 资源名
	URL           string `json:"url" validate:"required"`           // 资源路径
	RequestMethod string `json:"requestMethod" validate:"required"` // 请求方式
	ParentID      int    `json:"parentId" validate:"required"`      // 父资源id
	IsAnonymous   int    `json:"isAnonymous" validate:"required"`   // 是否匿名访问
}
