package dto

import "goBolg/model"

// ChatRecordDTO 代表聊天记录 DTO
type ChatRecordDTO struct {
	ChatRecordList model.ChatRecord `json:"chatRecordList"` // 聊天记录列表
	IPAddress      string           `json:"ipAddress"`      // IP 地址
	IPSource       string           `json:"ipSource"`       // IP 来源
}
