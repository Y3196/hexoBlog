package controller

import (
	"context"
	"goBolg/enums"
	"goBolg/exception"
	"goBolg/service"
	"goBolg/strategy/contxt"
	"goBolg/vo"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PhotoAlbumController struct {
	UploadStrategyContext *contxt.UploadStrategyContext
	PhotoAlbumService     service.PhotoAlbumService
}

func (p *PhotoAlbumController) SavePhotoAlbumCover(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Failed to get file"))
		return
	}

	log.Printf("Received file: %s", fileHeader.Filename)

	url, err := p.UploadStrategyContext.ExecuteUploadStrategy(fileHeader, enums.Photo) // 传递 pathEnum
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		c.JSON(http.StatusInternalServerError, vo.FailWithMessage("File upload failed"))
		return
	}

	log.Printf("File uploaded successfully: %s", url)
	c.JSON(http.StatusOK, vo.OkWithData(url))
}

func (p *PhotoAlbumController) SaveOrUpdatePhotoAlbum(c *gin.Context) {
	var photoAlbumVO vo.PhotoAlbumVO
	if err := c.ShouldBindJSON(&photoAlbumVO); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid input"))
		return
	}

	if err := p.PhotoAlbumService.SaveOrUpdatePhotoAlbum(context.Background(), photoAlbumVO); err != nil {
		log.Printf("Error saving or updating photo album: %v", err)
		if bizErr, ok := err.(*exception.BizError); ok {
			c.JSON(http.StatusInternalServerError, gin.H{"code": bizErr.Code, "message": bizErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	log.Println("Photo album saved or updated successfully")
	c.JSON(http.StatusOK, vo.OkWithMessage("插入或更新成功"))
}

// ListPhotoAlbumBacks 查看后台相册列表
// @Summary 查看后台相册列表
// @Description 查看后台相册列表
// @Tags 后台相册
// @Accept json
// @Produce json
// @Param condition query vo.ConditionVO true "条件"
// @Success 200 {object} Result{data=dto.PageResult{records=[]dto.PhotoAlbumBackDTO}} "相册列表"
// @Router /admin/photos/albums [get]
func (p *PhotoAlbumController) ListPhotoAlbumBacks(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.ShouldBindQuery(&condition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := p.PhotoAlbumService.ListPhotoAlbumBacks(c.Request.Context(), condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// ListPhotoAlbumBackInfos 获取相册信息列表
func (p *PhotoAlbumController) ListPhotoAlbumBackInfos(c *gin.Context) {
	ctx := c.Request.Context()
	photoAlbumDTOs, err := p.PhotoAlbumService.ListPhotoAlbumBackInfos(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(photoAlbumDTOs))
}

// GetPhotoAlbumBackByID 获取相册详细信息根据相册ID
func (p *PhotoAlbumController) GetPhotoAlbumBackByID(c *gin.Context) {
	albumID, err := strconv.Atoi(c.Param("albumId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	ctx := c.Request.Context()
	albumDTO, err := p.PhotoAlbumService.GetPhotoAlbumBackByID(ctx, albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(albumDTO))
}

// DeletePhotoAlbumByID 根据id删除相册
func (p *PhotoAlbumController) DeletePhotoAlbumByID(c *gin.Context) {
	albumID, err := strconv.Atoi(c.Param("albumId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithMessage("Invalid album ID"))
		return
	}

	ctx := c.Request.Context()
	if err := p.PhotoAlbumService.DeletePhotoAlbumByID(ctx, albumID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("相册删除成功"))
}

// ListPhotoAlbums 查看相册列表
func (p *PhotoAlbumController) ListPhotoAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	albums, err := p.PhotoAlbumService.ListPhotoAlbums(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(albums))
}
