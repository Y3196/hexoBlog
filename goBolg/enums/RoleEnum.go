package enums

// RoleEnum 角色枚举
type RoleEnum struct {
	RoleID int
	Name   string
	Label  string
}

var (
	// ADMIN 管理员
	ADMIN = RoleEnum{RoleID: 1, Name: "管理员", Label: "admin"}
	// USER 普通用户
	USER = RoleEnum{RoleID: 2, Name: "用户", Label: "user"}
	// TEST 测试账号
	TEST = RoleEnum{RoleID: 3, Name: "测试", Label: "test"}
)

// GetRoles 返回所有角色枚举值
func GetRoles() []RoleEnum {
	return []RoleEnum{ADMIN, USER, TEST}
}
