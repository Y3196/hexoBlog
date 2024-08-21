package controller

import (
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentController 处理评论相关的HTTP请求
type CommentController struct {
	CommentService service.CommentService
}

// ListComments 处理查询评论的HTTP请求
// @Summary List comments
// @Description Get all comments
// @Tags comments
// @Accept json
// @Produce json
// @Param topicId query int false "Topic ID"
// @Param type query int true "Comment Type"
// @Param current query int true "Current Page"
// @Param size query int true "Page Size"
// @Success 200 {object} vo.Response{data=vo.PageResult{recordList=[]dto.CommentDTO}}
// @Router /comments [get]
func (h *CommentController) ListComments(c *gin.Context) {
	topicID, _ := strconv.Atoi(c.Query("topicId"))
	commentType, _ := strconv.Atoi(c.Query("type"))
	current, _ := strconv.Atoi(c.Query("current"))
	size, _ := strconv.Atoi(c.Query("size"))

	if size == 0 {
		size = 10
	}

	commentVO := vo.CommentVO{
		TopicID: topicID,
		Type:    commentType,
	}

	comments, total, err := h.CommentService.ListComments(c.Request.Context(), commentVO, current, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve comments"))
		return
	}

	// Logging for debugging
	log.Printf("Comments fetched: %+v", comments)
	log.Printf("Total comments: %d", total)

	pageResult := vo.NewPageResult(comments, total)
	result := vo.OkWithData(pageResult)

	c.JSON(http.StatusOK, result)
}

// ListRepliesByCommentId 处理查询评论下回复的HTTP请求
// @Summary 查询评论下的回复
// @Description 获取指定评论的回复列表
// @Tags comments
// @Accept json
// @Produce json
// @Param commentId path int true "评论ID"
// @Success 200 {object} vo.Response{data=[]dto.ReplyDTO}
// @Router /comments/{commentId}/replies [get]
func (h *CommentController) ListRepliesByCommentId(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid comment ID"))
		return
	}

	// Log the comment ID
	log.Printf("Fetching replies for comment ID: %d", commentId)

	replies, err := h.CommentService.ListRepliesByCommentId(c.Request.Context(), commentId)
	if err != nil {
		log.Printf("Failed to retrieve replies for comment ID: %d, error: %v", commentId, err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve replies"))
		return
	}

	if len(replies) == 0 {
		log.Printf("No replies found for comment ID: %d", commentId)
	}

	log.Printf("Replies fetched for comment ID %d: %v", commentId, replies)
	result := vo.OkWithData(replies)
	c.JSON(http.StatusOK, result)
}

// SaveComment 添加评论
// @Summary 添加评论
// @Description 添加评论
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body vo.CommentVO true "评论信息"
// @Success 200 {object} vo.Response
// @Router /comments [post]
func (h *CommentController) SaveComment(c *gin.Context) {
	// 打印请求头以确认是否包含 Authorization 头部
	authHeader := c.GetHeader("Authorization")
	log.Printf("Authorization header: %s", authHeader)
	var commentVO vo.CommentVO
	// 解析请求的 JSON 数据
	if err := c.ShouldBindJSON(&commentVO); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request payload"))
		return
	}

	log.Printf("Received commentVO: %+v", commentVO)

	if err := vo.ValidateCommentVO(commentVO); err != nil {
		log.Printf("Validation failed: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Validation failed"))
		return
	}

	savedComment, err := h.CommentService.SaveComment(c.Request.Context(), commentVO)
	if err != nil {
		log.Printf("Failed to save comment: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to save comment"))
		return
	}

	log.Printf("Saved comment: %+v", savedComment)

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"recordList": []interface{}{
				savedComment,
			},
			"count": 1, // Adjust if necessary
		},
	}

	log.Printf("Returning comment data to front-end: %+v", response)
	c.JSON(http.StatusOK, response)
}

// SaveCommentLike 处理评论点赞的HTTP请求
// @Summary Save comment like
// @Description Like a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param commentId path int true "Comment ID"
// @Success 200 {object} vo.Response
// @Router /comments/{commentId}/like [post]
func (h *CommentController) SaveCommentLike(c *gin.Context) {
	// 获取评论ID
	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid comment ID"))
		return
	}

	// 调用CommentService的SaveCommentLike方法
	likeCount, err := h.CommentService.SaveCommentLike(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to save comment like"))
		return
	}

	// 返回成功响应，包含更新后的点赞数量
	response := vo.OkWithData(likeCount)
	c.JSON(http.StatusOK, response)
}

// UpdateCommentsReview 处理更新评论审核状态的HTTP请求
// @Summary Update comments review status
// @Description Update the review status of comments
// @Tags comments
// @Accept json
// @Produce json
// @Param reviewVO body vo.ReviewVO true "Review VO"
// @Success 200 {object} vo.Response{data=string}
// @Router /admin/comments/review [put]
func (h *CommentController) UpdateCommentsReview(c *gin.Context) {
	var reviewVO vo.ReviewVO
	if err := c.ShouldBindJSON(&reviewVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request data"))
		return
	}

	err := h.CommentService.UpdateCommentsReview(c.Request.Context(), reviewVO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to update comments review status"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("Comments review status updated successfully"))
}

// ListCommentBackDTO 处理查询后台评论的HTTP请求
// @Summary List back comments
// @Description Get all back comments
// @Tags comments
// @Accept json
// @Produce json
// @Param condition query vo.ConditionVO true "Condition VO"
// @Success 200 {object} vo.Response{data=vo.PageResult{recordList=[]dto.CommentBackDTO}}
// @Router /admin/comments [get]
func (h *CommentController) ListCommentBackDTO(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.ShouldBindQuery(&condition); err != nil {
		log.Printf("Invalid query parameters: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	// 设置默认分页参数
	ctx := c.Request.Context()
	ctx = utils.SetCurrentPage(ctx, &utils.Page{Current: 1, Size: 10})

	// 调用服务层方法获取评论数据
	pageResult, err := h.CommentService.ListCommentBackDTO(ctx, condition)
	if err != nil {
		log.Printf("Failed to retrieve comments: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve comments"))
		return
	}

	// 返回成功响应及评论数据
	c.JSON(http.StatusOK, vo.OkWithData(pageResult))
}

// DeleteComments 批量删除评论
// @Summary 删除评论
// @Description 批量删除评论
// @Tags 评论
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response{}
// @Router /admin/comments [delete]
func (h *CommentController) RemoveComments(c *gin.Context) {
	var commentIdList []int
	if err := c.ShouldBindJSON(&commentIdList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID list"})
		return
	}

	if err := h.CommentService.RemoveComments(c.Request.Context(), commentIdList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comments removed successfully"})
}
