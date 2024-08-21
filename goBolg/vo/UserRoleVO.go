package vo

// UserRoleVO 代表用户角色对象
type UserRoleVO struct {
	UserInfoId int    `json:"userInfoId" validate:"required"` // 用户id，不能为空
	Nickname   string `json:"nickname" validate:"required"`   // 用户昵称，不能为空
	RoleIdList []int  `json:"roleIdList" validate:"required"` // 角色id集合，不能为空
}
