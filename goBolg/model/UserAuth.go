package model

import (
	"fmt"
	"time"
)

// UserAuth 代表用户账号实体
type UserAuth struct {
	// 主键id
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 用户信息id
	UserInfoID int `json:"userInfoId" gorm:"column:user_info_id"`

	// 用户名
	Username string `json:"username" gorm:"column:username"`

	// 密码
	Password string `json:"password" gorm:"column:password"`

	// 登录类型
	LoginType int `json:"loginType" gorm:"column:login_type"`

	// 用户登录ip
	IPAddress string `json:"ipAddress" gorm:"column:ip_address"`

	// ip来源
	IPSource string `json:"ipSource" gorm:"column:ip_source"`

	// 创建时间
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`

	// 修改时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`

	// 最近登录时间
	LastLoginTime time.Time `gorm:"-"`
}

// String 返回结构体的字符串表示
func (ua UserAuth) String() string {
	return fmt.Sprintf("UserAuth{ID=%d, UserInfoID=%d, Username=%s, Password=%s, LoginType=%d, IPAddress=%s, IPSource=%s, CreateTime=%s, UpdateTime=%s, LastLoginTime=%s}",
		ua.ID, ua.UserInfoID, ua.Username, ua.Password, ua.LoginType, ua.IPAddress, ua.IPSource, ua.CreateTime, ua.UpdateTime, ua.LastLoginTime)
}
func (UserAuth) TableName() string {
	return "tb_user_auth" // 指定正确的表名
}
