package vo

import "github.com/go-playground/validator/v10"

// EmailVO 代表绑定邮箱
type EmailVO struct {
	// 邮箱
	Email string `json:"email" validate:"required,email"` // 邮箱地址，必填，且需为有效的邮箱格式

	// 验证码
	Code string `json:"code" validate:"required"` // 邮箱验证码，必填
}

// ValidateEmailVO 用于验证 EmailVO 结构体
func ValidateEmailVO(emailVO EmailVO) error {
	validate := validator.New()
	return validate.Struct(emailVO)
}
