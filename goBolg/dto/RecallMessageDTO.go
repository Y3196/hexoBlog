package dto

// RecallMessageDTO 代表撤回消息的 DTO
type RecallMessageDTO struct {
	ID      int  `json:"id"`      // 消息id
	IsVoice bool `json:"isVoice"` // 是否为语音
}
