package dto

import (
	"time"
)

// UserDetailDTO 代表用户信息 DTO
type UserDetailDTO struct {
	ID             int       `json:"id"`             // 用户账号id
	UserInfoID     int       `json:"userInfoId"`     // 用户信息id
	Email          string    `json:"email"`          // 邮箱号
	LoginType      int       `json:"loginType"`      // 登录方式
	Username       string    `json:"username"`       // 用户名
	Password       string    `json:"password"`       // 密码
	RoleList       []string  `json:"roleList"`       // 用户角色
	Nickname       string    `json:"nickname"`       // 用户昵称
	Avatar         string    `json:"avatar"`         // 用户头像
	Intro          string    `json:"intro"`          // 用户简介
	WebSite        string    `json:"webSite"`        // 个人网站
	ArticleLikeSet []string  `json:"articleLikeSet"` // 点赞文章集合
	CommentLikeSet []string  `json:"commentLikeSet"` // 点赞评论集合
	TalkLikeSet    []string  `json:"talkLikeSet"`    // 点赞说说集合
	IPAddress      string    `json:"ipAddress"`      // 用户登录ip
	IPSource       string    `json:"ipSource"`       // ip来源
	IsDisable      int       `json:"isDisable"`      // 是否禁用
	Browser        string    `json:"browser"`        // 浏览器
	OS             string    `json:"os"`             // 操作系统
	LastLoginTime  time.Time `json:"lastLoginTime"`  // 最近登录时间
}

// IsAccountNonExpired always returns true for the sake of simplicity
func (u *UserDetailDTO) IsAccountNonExpired() bool {
	return true
}

// IsAccountNonLocked checks if the account is disabled
func (u *UserDetailDTO) IsAccountNonLocked() bool {
	return u.IsDisable == 0
}

// IsCredentialsNonExpired always returns true for the sake of simplicity
func (u *UserDetailDTO) IsCredentialsNonExpired() bool {
	return true
}

// IsEnabled always returns true for the sake of simplicity
func (u *UserDetailDTO) IsEnabled() bool {
	return true
}

// GetAuthorities returns a list of authorities/roles
func (u *UserDetailDTO) GetAuthorities() []string {
	return u.RoleList
}
