package dto

// WeiboTokenDTO 代表微博 token 信息
type WeiboTokenDTO struct {
	UID         string `json:"uid"`          // 微博 uid
	AccessToken string `json:"access_token"` // 访问令牌
}
