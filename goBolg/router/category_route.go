package router

import (
	"github.com/gin-gonic/gin"
	"goBolg/app"
)

func SetupCategoryRoutes(router *gin.RouterGroup, controllers *app.Controllers, authMiddleWare gin.HandlerFunc) {

	router.GET("/articles", controllers.ArticleController.ListArticles)

	router.GET("/", controllers.BlogInfoController.GetBlogHomeInfo)

	router.GET("/about", controllers.BlogInfoController.GetAbout)

	router.GET("/articles/search", controllers.ArticleController.ListArticlesBySearch)

	router.POST("/report", controllers.BlogInfoController.Report)

	router.GET("/categories", controllers.CategoryController.ListCategories)

	router.GET("/links", controllers.FriendLinkController.ListFriendLinks)

	router.POST("/messages", controllers.MessageController.SaveMessage)

	router.GET("/messages", controllers.MessageController.ListMessages)

	router.GET("/photos/albums", controllers.PhotoAlbumController.ListPhotoAlbums)

	router.GET("/albums/:albumId/photos", controllers.PhotoController.ListPhotosByAlbumID)

	router.GET("/tags", controllers.TagController.ListTags)

	router.GET("/home/talks", controllers.TalkController.ListHomeTalks)

	router.GET("/talks/:talkId", controllers.TalkController.GetTalkById)

	router.GET("/talks", controllers.TalkController.ListTalks)

	router.GET("/comments", controllers.CommentController.ListComments)

	router.GET("/comments/:commentId/replies", controllers.CommentController.ListRepliesByCommentId)

	router.GET("/articles/:articleId", controllers.ArticleController.GetArticleById)

	router.POST("/talks/:talkId/like", authMiddleWare, controllers.TalkController.SaveTalkLike)

	router.PUT("/users/info", authMiddleWare, controllers.UserInfoController.UpdateUserInfo)

	router.POST("/users/avatar", authMiddleWare, controllers.UserInfoController.UpdateUserAvatar)

	router.POST("/users/email", authMiddleWare, controllers.UserInfoController.SaveUserEmail)

	router.GET("/users/code", controllers.UserAuthController.SendCode)

	router.POST("/register", controllers.UserAuthController.Register)

	router.PUT("/users/password", controllers.UserAuthController.UpdatePassword)

	router.GET("/admin/articles/:articleId", controllers.ArticleController.GetArticleBackById)

	router.POST("/admin/articles/images", controllers.ArticleController.SaveArticleImages)

	router.POST("/admin/categories", controllers.CategoryController.SaveOrUpdateCategory)

	router.GET("/admin/tags", controllers.TagController.ListTagBackDTO)

	router.DELETE("/admin/tags", controllers.TagController.DeleteTag)

	router.POST("/admin/tags", controllers.TagController.SaveOrUpdateTag)

	router.GET("/admin/tags/search", controllers.TagController.ListTagsBySearch)

	router.GET("/admin/categories", controllers.CategoryController.ListBackCategories)

	router.GET("/admin/categories/search", controllers.CategoryController.ListCategoriesBySearch)
}
