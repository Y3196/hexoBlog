package vo

// DeleteVO 代表逻辑删除
type DeleteVO struct {
	// ID 列表
	IDList   []int `json:"idList,omitempty"` // ID 列表
	IsDelete *int  `json:"isDelete"`         // 删除状态
}
