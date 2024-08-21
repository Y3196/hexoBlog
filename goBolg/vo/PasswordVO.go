package vo

// PasswordVO 代表密码修改的信息
type PasswordVO struct {
	OldPassword string `json:"oldPassword" validate:"required"`       // 旧密码
	NewPassword string `json:"newPassword" validate:"required,min=6"` // 新密码，最少6位
}
