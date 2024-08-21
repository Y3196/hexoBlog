package dto

// WeiboUserInfoDTO 代表微博用户信息
type WeiboUserInfoDTO struct {
	ScreenName string `json:"screen_name"` // 昵称
	AvatarHD   string `json:"avatar_hd"`   // 头像
}
