package dto

import "time"

// TalkBackDTO 代表后台说说
type TalkBackDTO struct {
	ID         int       `json:"id"`         // 说说id
	Nickname   string    `json:"nickname"`   // 昵称
	Avatar     string    `json:"avatar"`     // 头像
	Content    string    `json:"content"`    // 说说内容
	Images     string    `json:"images"`     // 图片
	ImgList    []string  `json:"imgList"`    // 图片列表
	IsTop      int       `json:"isTop"`      // 是否置顶
	Status     int       `json:"status"`     // 状态
	CreateTime time.Time `json:"createTime"` // 创建时间
}
