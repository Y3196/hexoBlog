package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"goBolg/app"
	"goBolg/config"
	"goBolg/handler"
	"goBolg/service"
)

func SetupRouter(app *app.Controllers, webSecurityConfig *config.WebSecurityConfig, userDetailsService service.UserDetailsService) *gin.Engine {
	router := gin.Default()
	webSecurityConfig.Configure(router)
	authMiddleware := handler.NewAuthenticationEntryPointImpl(userDetailsService).Middleware()

	// Swagger router
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	api := router.Group("")
	//api.Use(handler.SwaggerAuthMiddleware(), handler.PaginationMiddleware())

	// Setup normal routes
	SetupArticleRoutes(api, app, authMiddleware)
	/*SetupBlogInfoRoutes(api, app)*/
	SetupCategoryRoutes(api, app, authMiddleware)
	SetupCommentRoutes(api, app, authMiddleware)
	// Setup admin routes
	SetupAdminRoutes(api, app, authMiddleware)

	return router
}
