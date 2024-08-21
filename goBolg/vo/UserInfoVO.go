package vo

// UserInfoVO 代表用户信息对象
type UserInfoVO struct {
	Nickname string `json:"nickname" validate:"required"` // 用户昵称，不能为空
	Intro    string `json:"intro"`                        // 用户简介
	WebSite  string `json:"webSite"`                      // 个人网站
}
