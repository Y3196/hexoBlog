package vo

import (
	"time"
)

// VoiceVO 代表音频对象
type VoiceVO struct {
	// 消息类型
	Type int `json:"type" validate:"required"`

	// 文件内容
	File []byte `json:"file" validate:"required"` // 使用 []byte 表示文件内容

	// 用户id
	UserId int `json:"userId" validate:"required"`

	// 用户昵称
	Nickname string `json:"nickname" validate:"required"`

	// 用户头像
	Avatar string `json:"avatar" validate:"required"`

	// 聊天内容
	Content string `json:"content" validate:"required"`

	// 创建时间
	CreateTime time.Time `json:"createTime" validate:"required"`

	// 用户登录ip
	IpAddress string `json:"ipAddress" validate:"required"`

	// ip来源
	IpSource string `json:"ipSource" validate:"required"`
}
