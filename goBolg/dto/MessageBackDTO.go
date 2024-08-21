package dto

import (
	"time"
)

// MessageBackDTO 代表后台留言 DTO
type MessageBackDTO struct {
	ID             int       `json:"id"`             // 主键id
	IPAddress      string    `json:"ipAddress"`      // 用户ip
	IPSource       string    `json:"ipSource"`       // 用户ip地址
	Nickname       string    `json:"nickname"`       // 昵称
	Avatar         string    `json:"avatar"`         // 头像
	MessageContent string    `json:"messageContent"` // 留言内容
	IsReview       int       `json:"isReview"`       // 是否审核
	CreateTime     time.Time `json:"createTime"`     // 留言时间
}
