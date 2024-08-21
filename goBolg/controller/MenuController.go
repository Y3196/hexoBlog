package controller

import (
	"github.com/gin-gonic/gin"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"log"
	"net/http"
	"strconv"
)

type MenuController struct {
	MenuService service.MenuService
}

// ListMenus 查看菜单列表
// @Summary 查看菜单列表
// @Description 获取所有菜单
// @Tags admin
// @Accept json
// @Produce json
// @Param keywords query string false "Keywords"
// @Success 200 {object} vo.Response{data=[]dto.MenuDTO}
// @Router /admin/menus [get]
func (h *MenuController) ListMenus(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.ShouldBindQuery(&condition); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	result, err := h.MenuService.ListMenus(c.Request.Context(), condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve menus"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// SaveOrUpdateMenu 保存或更新菜单
// @Summary 新增或修改菜单
// @Description 新增或修改菜单
// @Tags admin
// @Accept json
// @Produce json
// @Param menuVO body vo.MenuVO true "菜单信息"
// @Success 200 {object} vo.Response{message=string}
// @Router /admin/menus [post]
func (h *MenuController) SaveOrUpdateMenu(c *gin.Context) {
	var menuVO vo.MenuVO
	if err := c.ShouldBindJSON(&menuVO); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request parameters"))
		return
	}

	// 打印日志调试
	log.Printf("Received MenuVO: %+v", menuVO)

	if err := h.MenuService.SaveOrUpdateMenu(c.Request.Context(), menuVO); err != nil {
		log.Printf("Error saving or updating menu: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to save or update menu"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("保存或更新菜单成功"))
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除菜单
// @Tags admin
// @Accept json
// @Produce json
// @Param menuId path int true "Menu ID"
// @Success 200 {object} vo.Response{data=string}
// @Failure 400 {object} vo.Response{data=string}
// @Failure 500 {object} vo.Response{data=string}
// @Router /admin/menus/{menuId} [delete]
func (h *MenuController) DeleteMenu(c *gin.Context) {
	idParam := c.Param("menuId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request parameters"))
		return
	}

	if err := h.MenuService.DeleteMenu(c.Request.Context(), uint(id)); err != nil {
		if err.Error() == "菜单下有角色关联" {
			c.JSON(http.StatusBadRequest, vo.FailWithMessage(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to delete menu"))
		}
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("删除成功"))
}

// ListMenuOptions 查看角色菜单选项
// @Summary 查看角色菜单选项
// @Description 获取所有角色菜单选项
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response{data=[]dto.LabelOptionDTO}
// @Router /admin/role/menus [get]
func (h *MenuController) ListMenuOptions(c *gin.Context) {
	result, err := h.MenuService.ListMenuOptions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve menu options"))
		return
	}
	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// ListUserMenus 查看当前用户菜单
// @Summary 查看当前用户菜单
// @Description 获取当前用户的所有菜单
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response{data=[]dto.UserMenuDTO}
// @Router /admin/user/menus [get]
func (h *MenuController) ListUserMenus(c *gin.Context) {
	user, ok := utils.GetLoginUser(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, vo.FailWithMessage("Failed to get login user"))
		return
	}

	result, err := h.MenuService.ListUserMenus(c.Request.Context(), uint(user.UserInfoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve user menus"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(result))
}
