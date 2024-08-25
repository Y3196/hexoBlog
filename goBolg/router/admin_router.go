package router

import (
	"github.com/gin-gonic/gin"
	"goBolg/app"
)

func SetupAdminRoutes(router *gin.RouterGroup, app *app.Controllers, authMiddleWare gin.HandlerFunc) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(authMiddleWare)
	{
		// 文章
		adminGroup.GET("/articles", app.ArticleController.ListArticleBacks)
		//	adminGroup.GET("/articles/:articleId", app.ArticleController.GetArticleBackById)
		adminGroup.POST("/articles", app.ArticleController.SaveOrUpdateArticle)
		//	adminGroup.POST("/articles/images", app.ArticleController.SaveArticleImages)
		adminGroup.PUT("/articles/top", app.ArticleController.UpdateArticleTop)
		adminGroup.PUT("/articles", app.ArticleController.UpdateArticleDelete)
		adminGroup.DELETE("/articles", app.ArticleController.DeleteArticles)
		// 网站
		adminGroup.GET("", app.BlogInfoController.GetBlogBackInfo)
		adminGroup.PUT("/website/config", app.BlogInfoController.UpdateWebsiteConfig)
		adminGroup.PUT("/about", app.BlogInfoController.UpdateAbout)
		// 分类
		//adminGroup.GET("/categories", app.CategoryController.ListBackCategories)
		//adminGroup.GET("/categories/search", app.CategoryController.ListCategoriesBySearch)
		//adminGroup.POST("/categories", app.CategoryController.SaveOrUpdateCategory)
		adminGroup.DELETE("/categories", app.CategoryController.DeleteCategories)
		//评论
		adminGroup.PUT("/comments/review", app.CommentController.UpdateCommentsReview)
		/*
			{
			  "type": 1,
			  "isReview": 1,
			  "keywords": "git"
			}
		*/
		adminGroup.GET("/comments", app.CommentController.ListCommentBackDTO)
		adminGroup.DELETE("/comments", app.CommentController.RemoveComments)

		// 友链
		adminGroup.GET("/links", app.FriendLinkController.ListFriendLinkDTO)
		adminGroup.POST("/links", app.FriendLinkController.SaveOrUpdateFriendLink)
		adminGroup.DELETE("/links", app.FriendLinkController.RemoveFriendLinks)

		//日志
		adminGroup.GET("/operation/logs", app.LogController.ListOperationLogs)

		//菜单
		adminGroup.GET("/menus", app.MenuController.ListMenus)
		/*
			{
			  "id": 220,
			  "name": "测试22",
			  "parent_id": null,
			  "order_num": 1,
			  "path": "/dashboard1",
			  "component": "dashboard/inde2x",
			  "icon": "dashboard",
			  "isHidden": 1,
			  "orderNum": 10
			}

		*/
		adminGroup.POST("/menus", app.MenuController.SaveOrUpdateMenu)
		adminGroup.DELETE("/menus/:menuId", app.MenuController.DeleteMenu)
		adminGroup.GET("/role/menus", app.MenuController.ListMenuOptions)
		adminGroup.GET("/user/menus", app.MenuController.ListUserMenus)

		//留言
		adminGroup.GET("/messages", app.MessageController.ListMessageBackDTO)
		adminGroup.PUT("/messages/review", app.MessageController.UpdateMessagesReview)
		adminGroup.DELETE("/messages", app.MessageController.DeleteMessages)

		//页面
		adminGroup.DELETE("/pages/:pageId", app.PageController.DeletePage)
		adminGroup.POST("/pages", app.PageController.SaveOrUpdatePage)
		adminGroup.GET("/pages", app.PageController.ListPages)

		//相册
		adminGroup.POST("/photos/albums/cover", app.PhotoAlbumController.SavePhotoAlbumCover)
		adminGroup.POST("/photos/albums", app.PhotoAlbumController.SaveOrUpdatePhotoAlbum)
		adminGroup.GET("/photos/albums", app.PhotoAlbumController.ListPhotoAlbumBacks)
		adminGroup.GET("/photos/albums/info", app.PhotoAlbumController.ListPhotoAlbumBackInfos)
		adminGroup.GET("/photos/albums/:albumId/info", app.PhotoAlbumController.GetPhotoAlbumBackByID)
		adminGroup.DELETE("/photos/albums/:albumId", app.PhotoAlbumController.DeletePhotoAlbumByID)

		//照片
		adminGroup.GET("/photos", app.PhotoController.ListPhotos)
		adminGroup.PUT("/photos", app.PhotoController.UpdatePhoto)
		adminGroup.POST("/photos", app.PhotoController.SavePhotos)
		adminGroup.PUT("/photos/album", app.PhotoController.UpdatePhotosAlbum) // 测试没有效果
		adminGroup.PUT("/photos/delete", app.PhotoController.UpdatePhotoDelete)
		adminGroup.DELETE("/photos", app.PhotoController.DeletePhotos)

		//资源
		adminGroup.POST("/resources", app.ResourceController.SaveOrUpdateResource)
		adminGroup.DELETE("/resources/:resourceId", app.ResourceController.DeleteResource)
		adminGroup.GET("/resources", app.ResourceController.ListResources)
		adminGroup.GET("/role/resources", app.ResourceController.ListResourceOption)

		//角色
		adminGroup.GET("/users/role", app.RoleController.ListUserRoles)
		adminGroup.GET("/roles", app.RoleController.ListRoles)
		adminGroup.DELETE("/roles", app.RoleController.DeleteRoles)
		adminGroup.POST("/role", app.RoleController.SaveOrUpdateRole) // 更新不知道更新的是什么？？？？？并且创建会进行填充更新时间

		//标签
		//	adminGroup.GET("/tags", app.TagController.ListTagBackDTO)
		//adminGroup.DELETE("/tags", app.TagController.DeleteTag)
		//adminGroup.POST("/tags", app.TagController.SaveOrUpdateTag) // 更新不知道更新的是什么？？？？？并且创建会进行填充更新时间
		//adminGroup.GET("/tags/search", app.TagController.ListTagsBySearch)

		// 说说
		adminGroup.POST("/talks/images", app.TalkController.SaveTalkImages)

		//后台
		adminGroup.GET("/users/area", app.UserAuthController.ListUserAreas)
	}
}
