package dto

// MessageDTO 代表留言 DTO
type MessageDTO struct {
	ID             int    `json:"id"`             // 主键id
	Nickname       string `json:"nickname"`       // 昵称
	Avatar         string `json:"avatar"`         // 头像
	MessageContent string `json:"messageContent"` // 留言内容
	Time           int    `json:"time"`           // 弹幕速度
}
