package dto

// SocialTokenDTO 代表社交登录的 Token 信息
type SocialTokenDTO struct {
	OpenID      string `json:"openId"`      // 开放id
	AccessToken string `json:"accessToken"` // 访问令牌
	LoginType   int    `json:"loginType"`   // 登录类型
}
