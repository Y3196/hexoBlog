package controller

import (
	"context"
	"goBolg/service"
	"goBolg/vo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResourceController struct {
	ResourceService service.ResourceService
}

// SaveOrUpdateResource 新增或修改资源
func (c *ResourceController) SaveOrUpdateResource(ctx *gin.Context) {
	var resourceVO vo.ResourceVO
	if err := ctx.ShouldBindJSON(&resourceVO); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request payload"))
		return
	}

	err := c.ResourceService.SaveOrUpdateResource(context.Background(), resourceVO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.FailWithMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, vo.Ok())
}

// DeleteResource 处理删除资源的 HTTP 请求
func (c *ResourceController) DeleteResource(ctx *gin.Context) {
	resourceId := ctx.Param("resourceId")

	// 将资源 ID 转换为 uint
	id, err := strconv.ParseUint(resourceId, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid resource ID"))
		return
	}

	err = c.ResourceService.DeleteResource(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, vo.OkWithMessage("删除成功"))
}

// ListResources 处理查看资源列表的 HTTP 请求
func (c *ResourceController) ListResources(ctx *gin.Context) {
	var conditionVO vo.ConditionVO
	if err := ctx.ShouldBindQuery(&conditionVO); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	resourceDTOList, err := c.ResourceService.ListResources(ctx, conditionVO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, vo.OkWithData(resourceDTOList))
}

// ListResourceOption 查看角色资源选项
func (c *ResourceController) ListResourceOption(ctx *gin.Context) {
	options, err := c.ResourceService.ListResourceOption(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, vo.OkWithData(options))
}
