package router

import (
	"github.com/gin-gonic/gin"
	"goBolg/app"
)

func SetupArticleRoutes(router *gin.RouterGroup, app *app.Controllers, userMiddleware gin.HandlerFunc) {
	articleGroup := router.Group("/articles")
	{
		articleGroup.GET("/archives", app.ArticleController.ListArchives)
		//articleGroup.GET("/:articleId", app.ArticleController.GetArticleById)
		// 传值 categoryId=68或者tagId=32
		articleGroup.GET("/condition", app.ArticleController.ListArticlesByCondition)
		articleGroup.POST("/:articleId/like", userMiddleware, app.ArticleController.SaveArticleLike)

	}
}
