package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"goBolg/enums"
	"goBolg/service"
	"goBolg/strategy/contxt"
	"goBolg/vo"
	"log"
	"net/http"
	"strconv"
)

type TalkController struct {
	TalkService           service.TalkService
	UploadStrategyContext *contxt.UploadStrategyContext
}

// ListHomeTalks 处理查看首页说说的请求
func (c *TalkController) ListHomeTalks(ctx *gin.Context) {
	talks, err := c.TalkService.ListHomeTalks(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, vo.OkWithData(talks))
}

// GetTalkById 根据ID获取说说
func (tc *TalkController) GetTalkById(c *gin.Context) {
	talkIdStr := c.Param("talkId")
	talkId, err := strconv.Atoi(talkIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.OkWithMessage("Invalid talk ID"))
		return
	}

	talkDTO, err := tc.TalkService.GetTalkById(c.Request.Context(), talkId)
	if err != nil {
		c.JSON(http.StatusNotFound, vo.FailWithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(talkDTO))
}

// ListTalks 获取所有说说
func (tc *TalkController) ListTalks(c *gin.Context) {
	pageResult, err := tc.TalkService.ListTalks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(pageResult))
}

// SaveTalkLike 点赞说说
func (tc *TalkController) SaveTalkLike(c *gin.Context) {
	talkIdStr := c.Param("talkId")
	talkId, err := strconv.Atoi(talkIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.OkWithMessage("Invalid talk ID"))
		return
	}

	likeCount, err := tc.TalkService.SaveTalkLike(c.Request.Context(), talkId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	// 返回最新的点赞量
	c.JSON(http.StatusOK, vo.OkWithData(likeCount))
}

func (c *TalkController) SaveTalkImages(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}

	log.Printf("Received file: %s", fileHeader.Filename)

	if c.UploadStrategyContext == nil {
		log.Printf("Error: UploadStrategyContext is nil")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Upload strategy context is not initialized"})
		return
	}

	if c.UploadStrategyContext.Config == nil {
		log.Printf("Error: UploadStrategyContext config is nil")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Upload strategy context config is not initialized"})
		return
	}

	// Log the upload mode to verify it's correctly set
	log.Printf("Upload mode: %s", c.UploadStrategyContext.Config.Upload.Mode)

	fileUrl, err := c.UploadStrategyContext.ExecuteUploadStrategy(fileHeader, enums.Talk)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
		return
	}

	log.Printf("File uploaded successfully: %s", fileUrl)
	ctx.JSON(http.StatusOK, gin.H{"fileUrl": fileUrl})
}
