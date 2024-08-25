package controller

import (
	"github.com/gin-gonic/gin"
	"goBolg/config"
	"goBolg/enums"
	"goBolg/service"
	"goBolg/strategy/contxt"
	"goBolg/vo"
	"log"
	"net/http"
	"strconv"
)

// ArticleController 文章控制器
type ArticleController struct {
	Service               service.ArticleService
	UploadStrategyContext *contxt.UploadStrategyContext
	//SearchStrategyContext *context.SearchStrategyContext
}

// ListArchives 查看文章归档
// @Summary List archives
// @Description Get all archives
// @Tags articles
// @Accept json
// @Produce json
// @Success 200 {object} vo.Result
// @Security BearerAuth
// @Router /archives [get]
func (controller *ArticleController) ListArchives(c *gin.Context) {
	log.Println("Controller: ListArchives called")

	// 获取分页参数
	current, err := strconv.Atoi(c.Query("current"))
	if err != nil || current < 1 {
		current = 1
	}

	archives, totalCount, err := controller.Service.ListArchives(current)
	if err != nil {
		log.Printf("Error in controller layer: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve archives"))
		return
	}

	log.Printf("Controller: ListArchives returned %d records", len(archives))
	c.JSON(http.StatusOK, vo.OkWithData(map[string]interface{}{
		"recordList": archives,
		"count":      totalCount,
	}))
}

// ListArticles 查看首页文章
// @Summary List articles
// @Description Get all articles
// @Tags articles
// @Accept json
// @Produce json
// @Success 200 {object} vo.Result
// @Security BearerAuth
// @Router /articles [get]
func (controller *ArticleController) ListArticles(c *gin.Context) {
	ctx := c.Request.Context()

	log.Println("Fetching list of articles") // 添加日志
	articles, err := controller.Service.ListArticles(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve articles"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(articles))
}

// ListArticleBacks 查看后台文章
// @Summary List back office articles
// @Description Get articles from back office with conditions
// @Tags admin
// @Accept json
// @Produce json
// @Param condition query string false "Search condition"
// @Success 200 {object} vo.Result
// @Security BearerAuth
// @Router /admin/articles [get]
func (controller *ArticleController) ListArticleBacks(c *gin.Context) {
	// 解析查询条件
	var conditionVO vo.ConditionVO
	if err := c.ShouldBindQuery(&conditionVO); err != nil {
		log.Printf("Error parsing conditions: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	// 调用服务层获取后台文章列表
	articlePageResult, err := controller.Service.ListArticleBacks(c.Request.Context(), conditionVO)
	if err != nil {
		log.Printf("Error retrieving back office articles: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve back office articles"))
		return
	}

	// 日志记录并返回结果
	c.JSON(http.StatusOK, vo.OkWithData(articlePageResult))
}

// GetArticleById 根据文章 ID 查看文章
// @Summary Get article by ID
// @Description Get article details by ID
// @Tags articles
// @Accept json
// @Produce json
// @Param articleId path int true "文章 ID"
// @Success 200 {object} vo.Result
// @Router /articles/{articleId} [get]
func (controller *ArticleController) GetArticleById(c *gin.Context) {
	articleIdStr := c.Param("articleId")
	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid article ID"))
		return
	}

	article, err := controller.Service.GetArticleById(c.Request.Context(), articleId)
	if err != nil {
		log.Printf("Error retrieving article: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve article"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(article))
}

// ListArticlesByCondition 根据条件查询文章
// @Summary 根据条件查询文章
// @Description 根据条件获取文章列表
// @Tags articles
// @Accept json
// @Produce json
// @Param condition query vo.ConditionVO false "查询条件"
// @Success 200 {object} vo.Result{data=dto.ArticlePreviewListDTO}
// @Security BearerAuth
// @Router /articles/condition [get]
// ArticleController - 控制器实现
// ListArticlesByCondition 处理根据条件查询文章的请求
func (controller *ArticleController) ListArticlesByCondition(c *gin.Context) {
	// 打印所有的查询参数
	for key, value := range c.Request.URL.Query() {
		log.Printf("Query Parameter: %s = %v", key, value)
	}

	// 解析查询条件
	var conditionVO vo.ConditionVO
	if err := c.ShouldBindQuery(&conditionVO); err != nil {
		log.Printf("Error parsing conditions: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	// 打印解析后的条件，检查是否正确解析
	log.Printf("ConditionVO before query: %+v", conditionVO)

	// 调用服务层获取文章列表
	articlePreviewListDTO, err := controller.Service.ListArticlesByCondition(c.Request.Context(), conditionVO)
	if err != nil {
		log.Printf("Error retrieving articles by condition: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve articles by condition"))
		return
	}

	// 日志记录并返回结果
	log.Printf("Controller: ListArticlesByCondition returned %d articles", len(articlePreviewListDTO.ArticlePreviewDTOList))
	c.JSON(http.StatusOK, vo.OkWithData(articlePreviewListDTO))
}

// GetArticleBackById 根据ID查看后台文章
// @Summary 根据id查看后台文章
// @Description 根据ID查看后台文章
// @Tags admin
// @Accept json
// @Produce json
// @Param articleId path int true "文章ID"
// @Success 200 {object} vo.Result{data=vo.ArticleVO}
// @Security BearerAuth
// @Router /admin/articles/{articleId} [get]
func (controller *ArticleController) GetArticleBackById(c *gin.Context) {
	articleIdStr := c.Param("articleId")
	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid article ID"))
		return
	}

	article, err := controller.Service.GetArticleBackById(c.Request.Context(), articleId)
	if err != nil {
		log.Printf("Error retrieving back office article: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve back office article"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(article))
}

// SaveArticleLike 点赞文章
// @Summary 点赞文章
// @Description 点赞文章
// @Tags articles
// @Accept json
// @Produce json
// @Param articleId path int true "文章id"
// @Success 200 {object} vo.Result
// @Router /articles/{articleId}/like [post]
func (controller *ArticleController) SaveArticleLike(c *gin.Context) {
	articleIdStr := c.Param("articleId")
	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid article ID"))
		return
	}

	// 调用服务层方法并获取返回的likeCount
	likeCount, err := controller.Service.SaveArticleLike(c.Request.Context(), articleId)
	if err != nil {
		log.Printf("Error processing like for article %d: %v", articleId, err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to process article like"))
		return
	}

	// 返回成功响应和likeCount
	c.JSON(http.StatusOK, vo.OkWithData(likeCount))
}

// SaveOrUpdateArticle 添加或修改文章
// @Summary Add or update an article
// @Description Save or update article details
// @Tags articles
// @Accept json
// @Produce json
// @Param article body vo.ArticleVO true "Article Information"
// @Success 200 {object} vo.Result
// @Security BearerAuth
// @Router /admin/articles [post]
func (controller *ArticleController) SaveOrUpdateArticle(c *gin.Context) {
	var articleVO vo.ArticleVO
	if err := c.ShouldBindJSON(&articleVO); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request payload"))
		return
	}

	err := controller.Service.SaveOrUpdateArticle(c.Request.Context(), articleVO)
	if err != nil {
		log.Printf("Error saving or updating article: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to save or update article"))
		return
	}

	c.JSON(http.StatusOK, vo.Ok())
}

// UpdateArticleTop 修改文章置顶状态
// @Summary 修改文章置顶
// @Description 修改文章置顶状态
// @Tags admin
// @Accept json
// @Produce json
// @Param articleTopVO body vo.ArticleTopVO true "文章置顶信息"
// @Success 200 {object} vo.Result
// @Router /admin/articles/top [put]
func (controller *ArticleController) UpdateArticleTop(c *gin.Context) {
	var articleTopVO vo.ArticleTopVO
	if err := c.ShouldBindJSON(&articleTopVO); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request payload"))
		return
	}

	log.Printf("Received ArticleTopVO: %+v", articleTopVO) // 打印接收到的VO

	if articleTopVO.IsTop == nil {
		log.Printf("IsTop is nil for ArticleTopVO: %+v", articleTopVO)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("IsTop cannot be nil"))
		return
	}

	if err := controller.Service.UpdateArticleTop(c.Request.Context(), articleTopVO); err != nil {
		log.Printf("Error updating article top status: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to update article top status"))
		return
	}

	c.JSON(http.StatusOK, vo.Ok())
}

// UpdateArticleDelete 恢复或删除文章
// @Summary 恢复或删除文章
// @Description 恢复或删除文章
// @Tags admin
// @Accept json
// @Produce json
// @Param deleteVO body vo.DeleteVO true "逻辑删除信息"
// @Success 200 {object} vo.Result
// @Router /admin/articles [put]
func (controller *ArticleController) UpdateArticleDelete(c *gin.Context) {
	var deleteVO vo.DeleteVO
	if err := c.ShouldBindJSON(&deleteVO); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request payload"))
		return
	}

	log.Printf("Received DeleteVO: %+v", deleteVO) // 打印接收到的VO

	if deleteVO.IDList == nil || len(deleteVO.IDList) == 0 {
		log.Printf("IDList is empty for DeleteVO: %+v", deleteVO)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("IDList cannot be empty"))
		return
	}

	if deleteVO.IsDelete == nil {
		log.Printf("IsDelete is nil for DeleteVO: %+v", deleteVO)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("IsDelete cannot be nil"))
		return
	}

	if err := controller.Service.UpdateArticleDelete(c.Request.Context(), deleteVO); err != nil {
		log.Printf("Error updating article delete status: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to update article delete status"))
		return
	}

	c.JSON(http.StatusOK, vo.Ok())
}

// DeleteArticles 删除文章
// @Summary 删除文章
// @Description 物理删除文章
// @Tags admin
// @Accept json
// @Produce json
// @Param articleIdList body []int true "文章 ID 列表"
// @Success 200 {object} vo.Result
// @Security BearerAuth
// @Router /admin/articles [delete]
func (controller *ArticleController) DeleteArticles(c *gin.Context) {
	var articleIdList []int
	if err := c.ShouldBindJSON(&articleIdList); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request payload"))
		return
	}

	if len(articleIdList) == 0 {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Article ID list cannot be empty"))
		return
	}

	err := controller.Service.DeleteArticles(c.Request.Context(), articleIdList)
	if err != nil {
		log.Printf("Error deleting articles: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to delete articles"))
		return
	}

	c.JSON(http.StatusOK, vo.Ok())
}

func (controller *ArticleController) ListArticlesBySearch(c *gin.Context) {
	var conditionVO vo.ConditionVO
	if err := c.ShouldBindQuery(&conditionVO); err != nil {
		log.Printf("Error parsing conditions: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request parameters"))
		return
	}

	ctx := c.Request.Context()

	articles, err := controller.Service.ListArticlesBySearch(ctx, conditionVO)
	if err != nil {
		log.Printf("Error searching articles: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to search articles"))
		return
	}
	c.JSON(http.StatusOK, vo.OkWithData(articles))
}

// SaveArticleImages handles the uploading of article images.
// @Summary 上传文章图片
// @Description 上传文章图片
// @Tags Articles
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文章图片"
// @Success 200 {object} vo.Result "图片上传成功后的URL"
// @Router /admin/articles/images [post]
// SaveArticleImages handles image upload requests from the frontend
func (ac *ArticleController) SaveArticleImages(c *gin.Context) {
	// 检查 UploadStrategyContext 是否为 nil
	if ac.UploadStrategyContext == nil {
		log.Printf("UploadStrategyContext is nil, reinitializing...")

		// 使用正确的配置文件路径
		configPath := "./config/application.yaml" // 修改路径指向你的 application.yaml 文件
		config, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config from %s: %v", configPath, err)
			c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Server configuration error"))
			return
		}

		ac.UploadStrategyContext, err = contxt.NewUploadStrategyContext(config)
		if err != nil {
			log.Fatalf("Failed to initialize UploadStrategyContext: %v", err)
			c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Server configuration error"))
			return
		}
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file from form: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Failed to get file"))
		return
	}

	log.Printf("Received file: %s", fileHeader.Filename)

	fileUrl, err := ac.UploadStrategyContext.ExecuteUploadStrategy(fileHeader, enums.Article)
	if err != nil {
		log.Printf("Recovered from panic in ExecuteUploadStrategy: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to upload file"))
		return
	}

	log.Printf("File uploaded successfully to: %s", fileUrl)
	c.JSON(http.StatusOK, vo.OkWithData(fileUrl))
}
