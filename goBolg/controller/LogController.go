package controller

import (
	"goBolg/service"
	"goBolg/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	OperationLogService service.OperationLogService
}

// ListOperationLogs 查看操作日志
// @Summary 查看操作日志
// @Description 获取所有操作日志
// @Tags admin
// @Accept json
// @Produce json
// @Param keywords query string false "Keywords"
// @Param current query int false "Current page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} vo.Response{data=vo.PageResult{recordList=[]dto.OperationLogDTO}}
// @Router /admin/operation/logs [get]
func (h *LogController) ListOperationLogs(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.ShouldBindQuery(&condition); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	result, err := h.OperationLogService.ListOperationLogs(c.Request.Context(), condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve operation logs"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// RemoveOperationLogs 删除操作日志
// @Summary 删除操作日志
// @Description 删除操作日志
// @Tags admin
// @Accept json
// @Produce json
// @Param ids body []uint true "Log IDs"
// @Success 200 {object} vo.Response{message=string}
// @Router /admin/operation/logs [delete]
func (h *LogController) RemoveOperationLogs(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request parameters"))
		return
	}

	if err := h.OperationLogService.RemoveOperationLogs(c.Request.Context(), ids); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to remove operation logs"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("删除成功"))
}
