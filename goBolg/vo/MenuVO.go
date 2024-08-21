package vo

// MenuVO 代表菜单的基本信息
type MenuVO struct {
	ID        uint   `json:"id,omitempty"`                  // 菜单id
	Name      string `json:"name" validate:"required"`      // 菜单名
	Icon      string `json:"icon" validate:"required"`      // 菜单icon
	Path      string `json:"path" validate:"required"`      // 路径
	Component string `json:"component" validate:"required"` // 组件
	OrderNum  int    `json:"orderNum" validate:"required"`  // 排序
	ParentId  *int   `json:"parentId,omitempty"`            // 父id
	IsHidden  *int   `json:"isHidden,omitempty"`            // 是否隐藏
}
