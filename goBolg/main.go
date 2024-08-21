package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"goBolg/app"
	"goBolg/config"
	"goBolg/handler"
	"goBolg/router"
	"goBolg/service/Impl"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @security BearerAuth
func main() {
	// 初始化应用
	application, err := app.Initialize()
	if err != nil {
		fmt.Printf("Initialization failed: %v\n", err)
		return
	}
	defer application.Close()

	// 初始化 UserDetailsService
	userDetailsService := Impl.NewUserDetailsServiceImpl(application.Controllers.UserAuthDao, application.Controllers.UserInfoDao, application.Controllers.RoleDao, application.Controllers.RedisService, nil)

	// 初始化 AuthHandler
	authHandler := handler.NewAuthHandler(userDetailsService)

	// 初始化 WebSecurityConfig
	webSecurityConfig := config.NewWebSecurityConfig(*authHandler)

	// 设置路由
	r := router.SetupRouter(application.Controllers, webSecurityConfig, userDetailsService)

	// 设置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8081"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	// 启动服务器
	serverPort := fmt.Sprintf(":%d", application.Config.Server.Port)
	err = r.Run(serverPort)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
}
