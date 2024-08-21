package controller

import (
	"goBolg/service"
	"goBolg/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	MessageService service.MessageService
}

// SaveMessage 添加留言
func (h *MessageController) SaveMessage(c *gin.Context) {
	var messageVO vo.MessageVO
	if err := c.ShouldBindJSON(&messageVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request parameters"))
		return
	}

	if err := h.MessageService.SaveMessage(c.Request.Context(), messageVO); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to save message"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("添加成功"))
}

// ListMessages 查询留言列表
func (h *MessageController) ListMessages(c *gin.Context) {
	messages, err := h.MessageService.ListMessages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(messages))
}

// ListMessageBackDTO 分页查询后台留言列表
func (h *MessageController) ListMessageBackDTO(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.ShouldBindQuery(&condition); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid query parameters"))
		return
	}

	result, err := h.MessageService.ListMessageBackDTO(c.Request.Context(), condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to retrieve messages"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// UpdateMessagesReview 更新留言审核状态
func (h *MessageController) UpdateMessagesReview(c *gin.Context) {
	var reviewVO vo.ReviewVO
	if err := c.ShouldBindJSON(&reviewVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request parameters"))
		return
	}

	if err := h.MessageService.UpdateMessagesReview(c.Request.Context(), reviewVO); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to update messages review"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("Review updated successfully"))
}

// DeleteMessages 删除留言
func (h *MessageController) DeleteMessages(c *gin.Context) {
	var messageIdList []uint
	if err := c.ShouldBindJSON(&messageIdList); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid request"))
		return
	}

	if err := h.MessageService.DeleteMessages(c.Request.Context(), messageIdList); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Failed to delete messages"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("Delete messages successfully"))
}
