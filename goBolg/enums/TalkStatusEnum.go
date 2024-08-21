package enums

// TalkStatusEnum 定义了说说的状态枚举
type TalkStatusEnum int

const (
	// Public 公开状态
	Public TalkStatusEnum = iota + 1
	// Secret 私密状态
	Secret
)

// String 返回 TalkStatusEnum 的描述
func (s TalkStatusEnum) String() string {
	switch s {
	case Public:
		return "公开"
	case Secret:
		return "私密"
	default:
		return "未知状态"
	}
}

// GetStatusDescription 根据 TalkStatusEnum 返回描述
func GetStatusDescription(status TalkStatusEnum) string {
	return status.String()
}
