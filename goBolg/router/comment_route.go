package router

import (
	"github.com/gin-gonic/gin"
	"goBolg/app"
)

func SetupCommentRoutes(router *gin.RouterGroup, controllers *app.Controllers, authMiddleWare gin.HandlerFunc) {
	commentGroup := router.Group("/comments")
	commentGroup.Use(authMiddleWare)
	{
		//	commentGroup.GET("", controllers.CommentController.ListComments)
		// /comments/732/replies
		//commentGroup.GET("/:commentId/replies", controllers.CommentController.ListRepliesByCommentId)
		/*
			{
			  "replyUserID": 1006,
			  "topicID": 68,
			  "commentContent": "thank",
			  "parentID": 732,
			  "type": 1
			}
		*/
		commentGroup.POST("", authMiddleWare, controllers.CommentController.SaveComment)
		commentGroup.POST("/:commentId/like", authMiddleWare, controllers.CommentController.SaveCommentLike)

	}
}
