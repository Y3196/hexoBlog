package vo

// WeiboLoginVO 代表微博登录信息
type WeiboLoginVO struct {
	// Code 是微博登录的授权码
	Code string `json:"code" validate:"required"` // 使用 validate 标签进行验证
}
