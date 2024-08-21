package dto

// FriendLinkDTO 代表友情链接 DTO
type FriendLinkDTO struct {
	ID          uint   `json:"id"`           // 链接 ID
	LinkName    string `json:"link_name"`    // 链接名
	LinkAvatar  string `json:"link_avatar"`  // 链接头像
	LinkAddress string `json:"link_address"` // 链接地址
	LinkIntro   string `json:"link_intro"`   // 介绍
}
