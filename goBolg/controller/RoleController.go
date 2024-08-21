package controller

import (
	"context"
	"goBolg/service"
	"goBolg/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct {
	RoleService service.RoleService
}

// ListUserRoles 查询用户角色选项
// @Summary 查询用户角色选项
// @Description 查询用户角色选项
// @Tags 角色
// @Produce json
// @Success 200 {object} utils.Result{data=[]dto.UserRoleDTO}
// @Router /admin/users/role [get]
func (controller *RoleController) ListUserRoles(c *gin.Context) {
	ctx := context.Background()
	userRoleDTOList, err := controller.RoleService.ListUserRoles(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Result{
			Code:    http.StatusInternalServerError,
			Message: "查询用户角色选项失败",
		})
		return
	}
	c.JSON(http.StatusOK, vo.Result{
		Code: http.StatusOK,
		Data: userRoleDTOList,
	})
}

// ListRoles 查询角色列表
func (controller *RoleController) ListRoles(c *gin.Context) {
	var conditionVO vo.ConditionVO
	if err := c.ShouldBindQuery(&conditionVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithData(err.Error()))
		return
	}

	// 调用服务层方法获取数据
	pageResult, err := controller.RoleService.ListRoles(c.Request.Context(), &conditionVO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(pageResult))
}

// DeleteRoles 删除角色
func (controller *RoleController) DeleteRoles(c *gin.Context) {
	var roleIdList []int
	if err := c.BindJSON(&roleIdList); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid input"))
		return
	}

	err := controller.RoleService.DeleteRoles(c, roleIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("删除成功"))
}

// SaveOrUpdateRole 保存或更新角色
func (controller *RoleController) SaveOrUpdateRole(c *gin.Context) {
	var roleVO vo.RoleVO

	if err := c.ShouldBindJSON(&roleVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithData(err.Error()))
		return
	}

	err := controller.RoleService.SaveOrUpdateRole(context.Background(), &roleVO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithDataAndMessage(err.Error(), "保存或更新角色失败"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("保存或更新角色成功"))
}
