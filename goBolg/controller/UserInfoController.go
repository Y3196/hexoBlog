package controller

import (
	"github.com/gin-gonic/gin"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"log"
	"net/http"
)

// UserController 结构体
type UserInfoController struct {
	UserInfoService service.UserInfoService
}

// UpdateUserInfo 更新用户信息的处理函数
func (c *UserInfoController) UpdateUserInfo(ctx *gin.Context) {
	// 解析请求体中的JSON数据
	var userInfoVO vo.UserInfoVO
	if err := ctx.ShouldBindJSON(&userInfoVO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 更新用户信息
	if err := c.UserInfoService.UpdateUserInfo(ctx.Request.Context(), &userInfoVO); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info: " + err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, vo.OkWithData("User info updated successfully"))
}

func (c *UserInfoController) UpdateUserAvatar(ctx *gin.Context) {
	// 获取上传的文件
	_, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// 使用获取到的fileHeader
	avatarURL, err := c.UserInfoService.UpdateUserAvatar(ctx.Request.Context(), fileHeader)
	if err != nil {
		log.Printf("Error updating avatar: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar: " + err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"avatar": avatarURL})
}

func (controller *UserInfoController) SaveUserEmail(c *gin.Context) {
	var emailVO vo.EmailVO
	if err := c.ShouldBindJSON(&emailVO); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Bound emailVO: %+v", emailVO)

	// 获取当前登录用户信息
	userDetailDTO, ok := utils.GetLoginUser(c.Request.Context())
	if !ok {
		log.Println("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}
	log.Printf("User loaded from context: %+v", userDetailDTO)
	log.Printf("Request context: %+v", c.Request.Context())
	err := controller.UserInfoService.SaveUserEmail(c.Request.Context(), emailVO)
	if err != nil {
		log.Printf("Error saving user email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
