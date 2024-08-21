package vo

// ReviewVO 代表审核信息
type ReviewVO struct {
	IDList   []int `json:"idList" validate:"required,dive,gt=0"` // id列表
	IsReview int   `json:"isReview" validate:"required"`         // 状态值
}
