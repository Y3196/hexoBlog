package controller

import (
	"goBolg/service"
	"goBolg/vo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FriendLinkController struct {
	FriendLinkService service.FriendLinkService
}

// ListFriendLinks 获取友链列表
// @Summary 获取友链列表
// @Description 获取所有友链
// @Tags 友链
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response{data=[]dto.FriendLinkDTO}
// @Failure 500 {object} vo.Response{message=string} "Failed to retrieve friend links"
// @Router /links [get]
func (h *FriendLinkController) ListFriendLinks(c *gin.Context) {
	friendLinks, err := h.FriendLinkService.ListFriendLinks(c.Request.Context())
	if err != nil {
		log.Printf("Failed to retrieve friend links: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve friend links"))
		return
	}
	log.Printf("Retrieved friend links: %+v", friendLinks)
	c.JSON(http.StatusOK, vo.OkWithData(friendLinks))
}

// ListFriendLinkDTO 查看后台友链列表
// @Summary 查看后台友链列表
// @Description 获取所有后台友链
// @Tags admin
// @Accept json
// @Produce json
// @Param keywords query string false "Keywords"
// @Success 200 {object} vo.Response{data=vo.PageResult{recordList=[]dto.FriendLinkBackDTO}}
// @Router /admin/links [get]
func (h *FriendLinkController) ListFriendLinkDTO(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.ShouldBindQuery(&condition); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	result, err := h.FriendLinkService.ListFriendLinkDTO(c.Request.Context(), condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve friend links"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// SaveOrUpdateFriendLink 保存或更新友链
// @Summary 保存或更新友链
// @Description 保存或更新友链信息
// @Tags admin
// @Accept json
// @Produce json
// @Param friendLink body vo.FriendLinkVO true "FriendLinkVO"
// @Success 200 {object} vo.Response{message=string} "保存或更新友链成功"
// @Failure 400 {object} vo.Response{message=string} "Invalid request parameters"
// @Failure 500 {object} vo.Response{message=string} "Failed to save friend link"
func (h *FriendLinkController) SaveOrUpdateFriendLink(c *gin.Context) {
	var friendLinkVO vo.FriendLinkVO
	if err := c.ShouldBindJSON(&friendLinkVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request parameters"))
		return
	}

	if err := h.FriendLinkService.SaveOrUpdateFriendLink(c.Request.Context(), friendLinkVO); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to save friend link"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("保存或更新友链成功"))
}

// RemoveFriendLinks 删除友链
// @Summary 删除友链
// @Description 删除友链信息
// @Tags admin
// @Accept json
// @Produce json
// @Param ids body []uint true "FriendLink IDs"
// @Success 200 {object} vo.Response{message=string} "删除成功"
// @Failure 400 {object} vo.Response{message=string} "Invalid request parameters"
// @Failure 500 {object} vo.Response{message=string} "Failed to remove friend links"
// @Router /admin/links [delete]
func (h *FriendLinkController) RemoveFriendLinks(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request parameters"))
		return
	}

	if err := h.FriendLinkService.RemoveFriendLinks(c.Request.Context(), ids); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to remove friend links"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("删除成功"))
}
