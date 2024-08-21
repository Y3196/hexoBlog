package dto

import "time"

// UserBackDTO 代表后台用户 DTO
type UserBackDTO struct {
	ID            int           `json:"id"`            // 用户id
	UserInfoID    int           `json:"userInfoId"`    // 用户信息id
	Avatar        string        `json:"avatar"`        // 头像
	Nickname      string        `json:"nickname"`      // 昵称
	RoleList      []UserRoleDTO `json:"roleList"`      // 用户角色
	LoginType     int           `json:"loginType"`     // 登录类型
	IPAddress     string        `json:"ipAddress"`     // 用户登录ip
	IPSource      string        `json:"ipSource"`      // ip来源
	CreateTime    time.Time     `json:"createTime"`    // 创建时间
	LastLoginTime time.Time     `json:"lastLoginTime"` // 最近登录时间
	IsDisable     int           `json:"isDisable"`     // 用户评论状态
	Status        int           `json:"status"`        // 状态
}
