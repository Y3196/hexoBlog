package dto

import "time"

// FriendLinkBackDTO 代表后台友情链接 DTO
type FriendLinkBackDTO struct {
	ID          int       `json:"id"`           // 链接 ID
	LinkName    string    `json:"link_name"`    // 链接名
	LinkAvatar  string    `json:"link_avatar"`  // 链接头像
	LinkAddress string    `json:"link_address"` // 链接地址
	LinkIntro   string    `json:"link_intro"`   // 链接介绍
	CreateTime  time.Time `json:"create_time"`  // 创建时间
}
