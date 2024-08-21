package controller

import (
	"goBolg/exception"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goBolg/service"
	"goBolg/vo"
)

type PhotoController struct {
	PhotoService service.PhotoService
}

// ListPhotos 获取后台照片列表
func (p *PhotoController) ListPhotos(c *gin.Context) {
	var condition vo.ConditionVO
	if err := c.ShouldBind(&condition); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithData(err.Error()))
		return
	}

	result, err := p.PhotoService.ListPhotos(c.Request.Context(), condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// UpdatePhoto 更新照片信息
func (p *PhotoController) UpdatePhoto(c *gin.Context) {
	var photoInfoVO vo.PhotoInfoVO
	if err := c.ShouldBindJSON(&photoInfoVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithData(err.Error()))
		return
	}

	if err := p.PhotoService.UpdatePhoto(c.Request.Context(), photoInfoVO); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("Photo updated successfully"))
}

// SavePhotos 保存照片信息
func (p *PhotoController) SavePhotos(c *gin.Context) {
	var photoVO vo.PhotoVO
	if err := c.ShouldBindJSON(&photoVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithData(err.Error()))
		return
	}

	if err := p.PhotoService.SavePhotos(c.Request.Context(), photoVO); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("Photos saved successfully"))
}

// 移动照片相册
func (p *PhotoController) UpdatePhotosAlbum(c *gin.Context) {
	var photoVO vo.PhotoVO
	if err := c.ShouldBindJSON(&photoVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithData(err.Error()))
		return
	}

	if err := p.PhotoService.UpdatePhotosAlbum(c.Request.Context(), photoVO); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("Photos album updated successfully"))
}

func (p *PhotoController) UpdatePhotoDelete(c *gin.Context) {
	var deleteVO vo.DeleteVO
	if err := c.ShouldBindJSON(&deleteVO); err != nil {
		c.JSON(http.StatusBadRequest, vo.FailWithData(err.Error()))
		return
	}

	if err := p.PhotoService.UpdatePhotoDelete(c.Request.Context(), deleteVO); err != nil {
		c.JSON(http.StatusInternalServerError, vo.FailWithDataAndMessage(err.Error(), "操作失败"))
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("Photos delete status updated successfully"))
}

// DeletePhotos 删除照片
func (p *PhotoController) DeletePhotos(c *gin.Context) {
	var photoIdList []int
	if err := c.ShouldBindJSON(&photoIdList); err != nil {
		vo.FailWithDataAndMessage(err.Error(), "参数错误")
		return
	}

	err := p.PhotoService.DeletePhotos(c.Request.Context(), photoIdList)
	if err != nil {
		vo.FailWithDataAndMessage(err.Error(), "操作失败")
		return
	}

	c.JSON(http.StatusOK, vo.OkWithMessage("删除成功"))
}

func (p *PhotoController) ListPhotosByAlbumID(c *gin.Context) {
	albumIDStr := c.Param("albumId")
	albumID, err := strconv.Atoi(albumIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid album ID"})
		return
	}

	// 获取分页参数
	current, size, err := getPagingParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid paging parameters"})
		return
	}

	result, err := p.PhotoService.ListPhotosByAlbumID(c.Request.Context(), albumID, current, size)
	if err != nil {
		if bizErr, ok := err.(*exception.BizError); ok {
			c.JSON(http.StatusInternalServerError, vo.FailWithDataAndMessage(bizErr.Code, bizErr.Message))
		} else {
			c.JSON(http.StatusInternalServerError, vo.FailWithMessage("Internal Server Error"))
		}
		return
	}

	c.JSON(http.StatusOK, vo.OkWithData(result))
}

// 获取分页参数
func getPagingParams(c *gin.Context) (int, int, error) {
	currentStr := c.Query("current")
	sizeStr := c.Query("size")

	current, err := strconv.Atoi(currentStr)
	if err != nil {
		return 0, 0, err
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return 0, 0, err
	}

	return current, size, nil
}
