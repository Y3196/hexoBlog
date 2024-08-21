package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"goBolg/constant"
	"goBolg/enums"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"log"
	"net/http"
	"strings"
)

type AuthenticationEntryPointImpl struct {
	UserDetailsService service.UserDetailsService
}

func NewAuthenticationEntryPointImpl(userDetailsService service.UserDetailsService) *AuthenticationEntryPointImpl {
	return &AuthenticationEntryPointImpl{
		UserDetailsService: userDetailsService,
	}
}

func (a *AuthenticationEntryPointImpl) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头部
		tokenString := c.GetHeader("Authorization")
		log.Printf("Authorization header: %s", tokenString)

		// 验证token
		if !a.isAuthenticated(c) {
			a.commence(c.Writer, c.Request)
			c.Abort()
			return
		}

		// 从上下文中获取userID
		userID, ok := c.Get("userID")
		if !ok || userID == nil {
			c.JSON(http.StatusUnauthorized, "userID is missing or invalid")
			return
		}

		userIDInt, ok := userID.(int)
		if !ok {
			log.Printf("userID is of unexpected type: %T", userID)
			c.JSON(http.StatusInternalServerError, "userID type assertion failed")
			return
		}
		log.Printf("userID extracted from context: %d", userIDInt)

		// 从数据库中加载用户详细信息
		user, err := a.UserDetailsService.LoadUserByID(c.Request.Context(), userIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Failed to load user")
			return
		}
		log.Printf("User loaded from LoadUserByID: %+v", user)

		// 将用户信息存入上下文
		ctx := context.WithValue(c.Request.Context(), constants.UserContextKey, user)
		c.Request = c.Request.WithContext(ctx)
		log.Printf("Set user in context: %+v", user)
		// 继续后续的处理
		c.Next()
	}
}

func (a *AuthenticationEntryPointImpl) isAuthenticated(c *gin.Context) bool {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		log.Println("Authorization header is missing")
		return false
	}

	// 移除 "Bearer " 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 验证JWT并获取声明信息
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		log.Printf("Failed to validate JWT: %v", err)
		return false
	}

	if claims == nil {
		log.Println("Claims are nil")
		return false
	}

	log.Printf("Authenticated claims: %+v", claims)

	// 将用户ID设置到Gin的上下文中
	c.Set("userID", claims.UserID)

	return true
}

func (a *AuthenticationEntryPointImpl) commence(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := vo.FailWithCodeAndMessage(enums.NO_LOGIN.Code, enums.NO_LOGIN.Desc)
	json.NewEncoder(w).Encode(response)
}
