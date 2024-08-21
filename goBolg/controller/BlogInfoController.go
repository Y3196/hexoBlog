package controller

import (
	"github.com/gin-gonic/gin"
	"goBolg/service"
	"goBolg/vo"
	"net/http"
)

// BlogInfoController 结构体用于存放博客信息服务
type BlogInfoController struct {
	BlogInfoService service.BlogInfoService
}

// GetBlogHomeInfo 处理GET请求，检索博客首页信息
func (ctrl *BlogInfoController) GetBlogHomeInfo(c *gin.Context) {
	info, err := ctrl.BlogInfoService.GetBlogHomeInfo(c.Request.Context())
	if err != nil {
		// 返回内部服务器错误和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 正常情况下返回状态码200和博客信息
	c.JSON(http.StatusOK, gin.H{"data": info})
}

// GetBlogBackInfo 处理 GET 请求，检索博客后台信息
func (ctrl *BlogInfoController) GetBlogBackInfo(c *gin.Context) {
	// 调用 rabbitService 层获取后台信息
	backInfo, err := ctrl.BlogInfoService.GetBlogBackInfo(c.Request.Context())
	if err != nil {
		// 如果有错误发生，返回 HTTP 500 错误码和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 如果没有错误，返回 HTTP 200 状态码和后台信息
	c.JSON(http.StatusOK, gin.H{"data": backInfo})
}

// UpdateWebsiteConfig 处理 PUT 请求，更新网站配置
func (ctrl *BlogInfoController) UpdateWebsiteConfig(c *gin.Context) {
	var websiteConfigVO vo.WebsiteConfigVO

	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&websiteConfigVO); err != nil {
		// 请求体格式错误，返回 HTTP 400 错误码和错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 调用服务层更新网站配置
	if err := ctrl.BlogInfoService.UpdateWebsiteConfig(c.Request.Context(), websiteConfigVO); err != nil {
		// 如果有错误发生，返回 HTTP 500 错误码和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功更新配置，返回 HTTP 200 状态码
	c.JSON(http.StatusOK, gin.H{"message": "Website configuration updated successfully"})
}

// GetAbout 处理 GET 请求，查看关于我信息
func (ctrl *BlogInfoController) GetAbout(c *gin.Context) {
	about, err := ctrl.BlogInfoService.GetAbout(c.Request.Context())
	if err != nil {
		// 如果有错误发生，返回 HTTP 500 错误码和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 正常情况下返回状态码200和关于我信息
	c.JSON(http.StatusOK, vo.OkWithData(about))
}

// UpdateAbout 处理 PUT 请求，更新关于我信息
func (ctrl *BlogInfoController) UpdateAbout(c *gin.Context) {
	var blogInfoVO vo.BlogInfoVO

	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&blogInfoVO); err != nil {
		// 请求体格式错误，返回 HTTP 400 错误码和错误信息
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request body"))
		return
	}

	// 调用服务层更新关于我信息
	if err := ctrl.BlogInfoService.UpdateAbout(c.Request.Context(), blogInfoVO); err != nil {
		// 如果有错误发生，返回 HTTP 500 错误码和错误信息
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	// 成功更新关于我信息，返回 HTTP 200 状态码
	c.JSON(http.StatusOK, vo.OkWithMessage("About information updated successfully"))
}

// Report 处理 POST 请求，上传访客信息
func (ctrl *BlogInfoController) Report(c *gin.Context) {
	ctrl.BlogInfoService.Report(c.Request.Context(), c.Request)

	// 返回成功响应
	c.JSON(http.StatusOK, vo.OkWithMessage("Visitor information reported successfully"))
}
