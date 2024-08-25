package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goBolg/service"
	"goBolg/vo"
	"log"
	"net/http"
)

// TagController 标签控制器
type TagController struct {
	Service service.TagService
}

// @Summary 查询标签列表
// @Description 获取所有标签的列表
// @Tags tags
// @Produce json
// @Success 200 {object} vo.PageResult[dto.TagDTO]
// @Router /tags [get]
func (tc *TagController) ListTags(c *gin.Context) {
	pageResult, err := tc.Service.ListTags(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}
	c.JSON(http.StatusOK, vo.OkWithData(pageResult))
}

// @Summary 查询后台标签列表
// @Description 获取后台标签的分页列表
// @Tags Tags
// @Accept json
// @Produce json
// @Param condition query vo.ConditionVO true "查询条件"
// @Success 200 {object} vo.PageResult[vo.TagBackDTO] "成功返回标签列表"
// @Failure 500 {object} map[string]string "内部服务器错误"
// @Router /admin/tags [get]
func (tc *TagController) ListTagBackDTO(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.ShouldBindQuery(&condition); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithData(err.Error()))
		return
	}

	pageResult, err := tc.Service.ListTagBackDTO(c, condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(pageResult))
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Description 根据标签ID列表删除标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param tagIdList body []int true "标签ID列表"
// @Success 200 {object} utils.Result
// @Failure 400 {object} utils.Result
// @Router /admin/tags [delete]
func (c *TagController) DeleteTag(ctx *gin.Context) {
	var tagIdList []int
	if err := ctx.ShouldBindJSON(&tagIdList); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Result{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Failed to bind JSON: %v", err),
		})
		return
	}

	log.Printf("Received tag ID list: %v", tagIdList) // Debug log to see received data

	err := c.Service.DeleteTag(ctx.Request.Context(), tagIdList)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Result{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, vo.Result{
		Code:    http.StatusOK,
		Message: "删除成功",
	})
}

// SaveOrUpdateTag 添加或修改标签
// @Summary 添加或修改标签
// @Description 添加或修改标签信息
// @Tags 标签
// @Accept json
// @Produce json
// @Param tagVO body vo.TagVO true "标签信息"
// @Success 200 {object} utils.Result
// @Failure 400 {object} utils.Result
// @Router /admin/tags [post]
func (c *TagController) SaveOrUpdateTag(ctx *gin.Context) {
	var tagVO vo.TagVO
	if err := ctx.ShouldBindJSON(&tagVO); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Result{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	tag, err := c.Service.SaveOrUpdateTag(ctx.Request.Context(), tagVO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Result{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, vo.Result{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data:    tag,  // 返回更新后的标签数据
		Flag:    true, // 确保这个值为 true
	})
}

func (ctl *TagController) ListTagsBySearch(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.BindQuery(&condition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tags, err := ctl.Service.ListTagsBySearch(c, condition)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.Result{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(tags))
}
