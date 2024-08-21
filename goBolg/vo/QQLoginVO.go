package vo

// QQLoginVO 代表QQ登录信息
type QQLoginVO struct {
	OpenID      string `json:"openId" validate:"required"`      // QQ openId
	AccessToken string `json:"accessToken" validate:"required"` // QQ accessToken
}
