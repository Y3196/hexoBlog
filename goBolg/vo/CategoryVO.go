package vo

import (
	"github.com/go-playground/validator/v10"
)

// CategoryVO 代表分类信息
type CategoryVO struct {
	ID           int    `json:"id,omitempty"`                     // 分类id
	CategoryName string `json:"categoryName" validate:"required"` // 分类名，必填
}

// ValidateCategoryVO 验证 CategoryVO 实例  ValidateCategoryVO: 这是一个辅助函数，使用 validator 库验证 CategoryVO 实例的必填字段
func ValidateCategoryVO(category CategoryVO) error {
	validate := validator.New()
	return validate.Struct(category)
}
