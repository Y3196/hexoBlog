package enums

// LoginTypeEnum 登录方式枚举
type LoginTypeEnum struct {
	Type     int
	Desc     string
	Strategy string
}

var (
	// EMAIL 邮箱登录
	EMAIL = LoginTypeEnum{Type: 1, Desc: "邮箱登录", Strategy: ""}
	// QQ QQ登录
	QQ = LoginTypeEnum{Type: 2, Desc: "QQ登录", Strategy: "qqLoginStrategyImpl"}
	// WEIBO 微博登录
	WEIBO = LoginTypeEnum{Type: 3, Desc: "微博登录", Strategy: "weiboLoginStrategyImpl"}
)

// GetLoginTypes 返回所有登录方式枚举值
func GetLoginTypes() []LoginTypeEnum {
	return []LoginTypeEnum{EMAIL, QQ, WEIBO}
}
