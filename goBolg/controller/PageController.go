package controller

import (
	"github.com/gin-gonic/gin"
	"goBolg/service"
	"goBolg/vo"
	"net/http"
	"strconv"
)

type PageController struct {
	PageService service.PageService
}

// DeletePage 删除页面
func (h *PageController) DeletePage(c *gin.Context) {
	pageID, err := strconv.Atoi(c.Param("pageId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid page ID"))
		return
	}

	if err := h.PageService.DeletePage(c.Request.Context(), uint(pageID)); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to delete page"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("Delete page successfully"))
}

// SaveOrUpdatePage 保存或更新页面
func (h *PageController) SaveOrUpdatePage(c *gin.Context) {
	var pageVO vo.PageVO
	if err := c.ShouldBindJSON(&pageVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request parameters"))
		return
	}

	if err := h.PageService.SaveOrUpdatePage(c.Request.Context(), pageVO); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to save or update page"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("Page saved or updated successfully"))
}

// ListPages 获取页面列表
func (h *PageController) ListPages(c *gin.Context) {
	pages, err := h.PageService.ListPages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve pages"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(pages))
}
