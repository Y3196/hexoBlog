package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"goBolg/service"
	"goBolg/vo"
	"net/http"
)

// CategoryController 结构体用于存放分类服务
type CategoryController struct {
	CategoryService service.CategoryService
}

// ListCategories 处理 GET 请求，查看分类列表
// @Summary 查看分类列表
// @Description 获取所有分类及其文章数量
// @Tags 分类模块
// @Produce json
// @Success 200 {object} vo.PageResult[dto.CategoryDTO]
// @Failure 500 {object} vo.Response{error=string}
// @Router /categories [get]
func (ctrl *CategoryController) ListCategories(c *gin.Context) {
	// 调用服务层获取分类列表
	result := ctrl.CategoryService.ListCategories(c.Request.Context())

	// 返回状态码200和分类列表
	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// ListBackCategories 处理 GET 请求，查看后台分类列表
// @Summary 查看后台分类列表
// @Description 获取后台分类列表，支持分页
// @Tags 分类模块
// @Produce json
// @Param condition query vo.ConditionVO true "查询条件"
// @Success 200 {object} vo.Result{data=vo.PageResult[dto.CategoryBackDTO]}
// @Failure 500 {object} vo.Result{message=string}
// @Router /admin/categories [get]
func (ctrl *CategoryController) ListBackCategories(c *gin.Context) {
	// 解析查询参数到 ConditionVO
	var condition vo.ConditionVO
	if err := c.ShouldBindQuery(&condition); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	// 调用服务层获取后台分类列表
	result, err := ctrl.CategoryService.ListBackCategories(c.Request.Context(), condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to get categories"))
		return
	}

	// 返回状态码200和分类列表
	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// ListCategoriesBySearch 处理 GET 请求，搜索文章分类
// @Summary 搜索文章分类
// @Description 根据搜索条件获取分类列表
// @Tags 分类模块
// @Produce json
// @Param keywords query string false "搜索关键词"
// @Success 200 {object} vo.Result{data=vo.PageResult[dto.CategoryOptionDTO]}
// @Failure 500 {object} vo.Result{message=string}
// @Router /admin/categories/search [get]
func (ctrl *CategoryController) ListCategoriesBySearch(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.ShouldBindQuery(&condition); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	result, err := ctrl.CategoryService.ListCategoriesBySearch(c.Request.Context(), condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to get categories"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// SaveOrUpdateCategory 处理 POST 请求，保存或更新分类
// @Summary 保存或更新分类
// @Description 保存或更新分类信息
// @Tags 分类模块
// @Param categoryVO body vo.CategoryVO true "分类信息"
// @Success 200 {object} vo.Result{data=string}
// @Failure 400 {object} vo.Result{message=string}
// @Failure 500 {object} vo.Result{message=string}
// @Router /admin/categories [post]
func (ctrl *CategoryController) SaveOrUpdateCategory(c *gin.Context) {
	var categoryVO vo.CategoryVO
	if err := c.ShouldBindJSON(&categoryVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request payload"))
		return
	}

	if err := ctrl.CategoryService.SaveOrUpdateCategory(c.Request.Context(), categoryVO); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("分类保存或更新成功"))
}

// DeleteCategories 处理 DELETE 请求，删除分类
// @Summary 删除分类
// @Description 删除分类
// @Tags 分类模块
// @Produce json
// @Param categoryIDList body []int true "分类ID列表"
// @Success 200 {object} vo.Result{data=string}
// @Failure 400 {object} vo.Result{message=string}
// @Failure 500 {object} vo.Result{message=string}
// @Router /admin/categories [delete]
func (ctrl *CategoryController) DeleteCategories(c *gin.Context) {
	var categoryIDList []int
	if err := c.ShouldBindJSON(&categoryIDList); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request payload"))
		return
	}

	if err := ctrl.CategoryService.DeleteCategory(context.Background(), categoryIDList); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("分类删除成功"))
}
