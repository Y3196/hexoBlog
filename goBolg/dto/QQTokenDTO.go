package dto

// QQTokenDTO 代表 QQ Token 信息
type QQTokenDTO struct {
	OpenID   string `json:"openid"`    // openid
	ClientID string `json:"client_id"` // 客户端id
}
