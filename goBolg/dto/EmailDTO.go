package dto

// EmailDTO 代表邮件 DTO
type EmailDTO struct {
	Email   string `json:"email"`   // 邮箱号
	Subject string `json:"subject"` // 主题
	Content string `json:"content"` // 内容
}
