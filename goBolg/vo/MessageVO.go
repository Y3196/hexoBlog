package vo

// MessageVO 代表留言的基本信息
type MessageVO struct {
	Nickname       string `json:"nickname" validate:"required"`       // 昵称
	Avatar         string `json:"avatar" validate:"required"`         // 头像
	MessageContent string `json:"messageContent" validate:"required"` // 留言内容
	Time           int    `json:"time" validate:"required"`           // 弹幕速度
}
