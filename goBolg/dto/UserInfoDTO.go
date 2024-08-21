package dto

import (
	"time"
)

// UserInfoDTO 代表用户信息 DTO
type UserInfoDTO struct {
	ID             int                      `json:"id"`             // 用户账号id
	UserInfoID     int                      `json:"userInfoId"`     // 用户信息id
	Email          string                   `json:"email"`          // 邮箱号
	LoginType      int                      `json:"loginType"`      // 登录方式
	Username       string                   `json:"username"`       // 用户名
	Nickname       string                   `json:"nickname"`       // 用户昵称
	Avatar         string                   `json:"avatar"`         // 用户头像
	Intro          string                   `json:"intro"`          // 用户简介
	WebSite        string                   `json:"webSite"`        // 个人网站
	ArticleLikeSet map[interface{}]struct{} `json:"articleLikeSet"` // 点赞文章集合
	CommentLikeSet map[interface{}]struct{} `json:"commentLikeSet"` // 点赞评论集合
	TalkLikeSet    map[interface{}]struct{} `json:"talkLikeSet"`    // 点赞说说集合
	IPAddress      string                   `json:"ipAddress"`      // 用户登录ip
	IPSource       string                   `json:"ipSource"`       // ip来源
	LastLoginTime  time.Time                `json:"lastLoginTime"`  // 最近登录时间
}
