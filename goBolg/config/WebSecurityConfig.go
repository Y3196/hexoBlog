package config

import (
	"github.com/gin-gonic/gin"

	"goBolg/handler"
)

// WebSecurityConfig holds the security configuration.
type WebSecurityConfig struct {
	AuthHandler handler.AuthHandler
}

// NewWebSecurityConfig creates a new WebSecurityConfig instance.
func NewWebSecurityConfig(authHandler handler.AuthHandler) *WebSecurityConfig {
	return &WebSecurityConfig{
		AuthHandler: authHandler,
	}
}

// Configure sets up the HTTP routes and middleware.
func (w *WebSecurityConfig) Configure(router *gin.Engine) {
	router.POST("/login", w.AuthHandler.Login)
	router.GET("/logout", w.AuthHandler.Logout)
}
