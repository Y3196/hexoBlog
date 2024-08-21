package vo

// UserVO 代表用户账号对象
type UserVO struct {
	Username string `json:"username" validate:"required,email"` // 用户名，不能为空且必须是有效的邮箱
	Password string `json:"password" validate:"required,min=6"` // 密码，不能为空且长度不能少于6位
	Code     string `json:"code" validate:"required"`           // 验证码，不能为空
}
