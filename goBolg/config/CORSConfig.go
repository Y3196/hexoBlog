// config/cors_config.go
package config

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 返回一个 CORS 中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.GetHeader("Origin")

		// 必须，接受指定域的请求，可以使用 * 但不安全
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 必须，设置服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")

		// 服务器支持的所有头信息字段
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token, Authorization")

		// 可选，设置 XMLHttpRequest 的响应对象能拿到的额外字段
		c.Header("Access-Control-Expose-Headers", "Authorization")

		// 可选，是否允许后续请求携带认证信息 Cookie
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有 OPTIONS 方法
		if method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
