package dto

// WebsocketMessageDTO 代表 WebSocket 消息
type WebsocketMessageDTO struct {
	Type int         `json:"type"` // 类型
	Data interface{} `json:"data"` // 数据
}
