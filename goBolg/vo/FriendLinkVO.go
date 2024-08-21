package vo

import "github.com/go-playground/validator/v10"

// FriendLinkVO 代表友链
type FriendLinkVO struct {
	// 友链id
	ID uint `json:"id"`

	// 链接名
	LinkName string `json:"linkName" validate:"required"` // 链接名，必填

	// 链接头像
	LinkAvatar string `json:"linkAvatar" validate:"required"` // 链接头像，必填

	// 链接地址
	LinkAddress string `json:"linkAddress" validate:"required"` // 链接地址，必填

	// 介绍
	LinkIntro string `json:"linkIntro" validate:"required"` // 链接介绍，必填
}

// ValidateFriendLinkVO 用于验证 FriendLinkVO 结构体
func ValidateFriendLinkVO(friendLinkVO FriendLinkVO) error {
	validate := validator.New()
	return validate.Struct(friendLinkVO)
}
