package vo

// UserDisableVO 代表用户禁用状态
type UserDisableVO struct {
	ID        int `json:"id" validate:"required"`        // 用户id，不能为空
	IsDisable int `json:"isDisable" validate:"required"` // 置顶状态，不能为空
}
