package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"goBolg/dto"
	"goBolg/service"
	"goBolg/utils"

	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication related requests.
type AuthHandler struct {
	UserDetailsService service.UserDetailsService
	SessionStore       *sessions.CookieStore
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(userDetailsService service.UserDetailsService) *AuthHandler {
	return &AuthHandler{
		UserDetailsService: userDetailsService,
	}
}

// Login handles user login requests.
func (h *AuthHandler) Login(c *gin.Context) {
	var loginDTO dto.LoginDTO

	// 从 URL 查询参数中获取用户名和密码
	loginDTO.Username = c.PostForm("username")
	loginDTO.Password = c.PostForm("password")

	if loginDTO.Username == "" || loginDTO.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码不能为空"})
		return
	}
	// 获取用户信息
	userDetail, err := h.UserDetailsService.LoadUserByUsername(c.Request, c.Request.Context(), loginDTO.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 比对密码
	if err := bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(loginDTO.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 生成 JWT
	token, err := utils.GenerateJWT(userDetail.ID, userDetail.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 在返回响应之前，先将用户ID存储到上下文中
	c.Set("userID", userDetail.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user":    userDetail,
		"token":   token,
	})
	c.Set("userID", userDetail.ID)
	// 更新用户信息

	go h.UserDetailsService.UpdateUserInfo(c.Request, userDetail)
}

// Logout handles the user logout request.
func (h *AuthHandler) Logout(c *gin.Context) {
	// 使 JWT 失效的逻辑（如果需要）可以在这里处理

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
