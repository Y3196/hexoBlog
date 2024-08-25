package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"goBolg/service"
	"goBolg/vo"
	"log"
	"net/http"
	"strconv"
)

type UserAuthController struct {
	UserAuthService service.UserAuthService
}

func (controller *UserAuthController) SendCode(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	err := controller.UserAuthService.SendCode(username)
	if err != nil {
		log.Printf("Error sending code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"flag": true, "message": "验证码已发送"})
}

// Register handles user registration.
func (controller *UserAuthController) Register(c *gin.Context) {
	var user vo.UserVO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"flag": false, "message": err.Error()})
		return
	}

	if err := controller.UserAuthService.Register(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"flag": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"flag": true, "message": "注册成功"})
}

func (controller *UserAuthController) UpdatePassword(c *gin.Context) {
	var user vo.UserVO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"flag": false, "message": err.Error()})
		return
	}

	err := controller.UserAuthService.UpdatePassword(c.Request.Context(), &user)
	if err != nil {
		log.Printf("Error updating password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"flag": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"flag": true, "message": "密码更新成功"})
}

func (controller *UserAuthController) ListUserAreas(c *gin.Context) {
	var condition vo.ConditionVO

	typeStr := c.Query("type")
	if typeStr != "" {
		typeInt, err := strconv.Atoi(typeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"flag": false, "message": "Invalid type parameter"})
			return
		}
		condition.Type = &typeInt
	}

	userAreas, err := controller.UserAuthService.ListUserAreas(context.Background(), condition)
	if err != nil {
		log.Printf("Error listing user areas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"flag": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"flag": true, "data": userAreas})
}
