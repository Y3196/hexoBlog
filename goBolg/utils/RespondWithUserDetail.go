package utils

import (
	"github.com/gin-gonic/gin"
	"goBolg/dto"
	"net/http"
)

// RespondWithUserDetail 封装用户详细信息的成功响应
func RespondWithUserDetail(c *gin.Context, message string, userDetail *dto.UserDetailDTO) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data": gin.H{
			"id":             userDetail.ID,
			"username":       userDetail.Username,
			"userInfoId":     userDetail.UserInfoID,
			"email":          userDetail.Email,
			"avatar":         userDetail.Avatar,
			"nickname":       userDetail.Nickname,
			"intro":          userDetail.Intro,
			"articleLikeSet": userDetail.ArticleLikeSet,
			"commentLikeSet": userDetail.CommentLikeSet,
			"talkLikeSet":    userDetail.TalkLikeSet,
			"ipAddress":      userDetail.IPAddress,
			"ipSource":       userDetail.IPSource,
			"lastLoginTime":  userDetail.LastLoginTime.UnixMilli(), // 时间转换为毫秒
			"loginType":      userDetail.LoginType,
		},
	})
}
