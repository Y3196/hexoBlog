package dto

import (
	"strconv"
	"strings"
	"time"
)

// RoleDTO 代表角色的 DTO
type RoleDTO struct {
	ID             int       `json:"id"`             // 角色id
	RoleName       string    `json:"roleName"`       // 角色名
	RoleLabel      string    `json:"roleLabel"`      // 角色标签
	CreateTime     time.Time `json:"createTime"`     // 创建时间
	IsDisable      int       `json:"isDisable"`      // 是否禁用
	ResourceIDList []int     `json:"resourceIdList"` // 资源id列表
	MenuIDList     []int     `json:"menuIdList"`     // 菜单id列表
}

// parseConcatList 解析拼接的ID列表
func ParseConcatList(concatStr string) []int {
	if concatStr == "" {
		return nil
	}
	var result []int
	ids := strings.Split(concatStr, ",")
	for _, id := range ids {
		if parsedID, err := strconv.Atoi(id); err == nil {
			result = append(result, parsedID)
		}
	}
	return result
}
