package dto

import "time"

// UserOnlineDTO 代表在线用户
type UserOnlineDTO struct {
	UserInfoID    int       `json:"userInfoId"`    // 用户信息id
	Nickname      string    `json:"nickname"`      // 用户昵称
	Avatar        string    `json:"avatar"`        // 用户头像
	IPAddress     string    `json:"ipAddress"`     // 用户登录ip
	IPSource      string    `json:"ipSource"`      // ip来源
	Browser       string    `json:"browser"`       // 浏览器
	OS            string    `json:"os"`            // 操作系统
	LastLoginTime time.Time `json:"lastLoginTime"` // 最近登录时间
}
