package Impl

import (
	"context"
	"fmt"
	"goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/enums"
	"goBolg/exception"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
)

type UserDetailsServiceImpl struct {
	userAuthDao  dao.UserAuthDao
	userInfoDao  dao.UserInfoDao
	roleDao      dao.RoleDao
	redisService service.RedisService
	store        *sessions.CookieStore
}

func NewUserDetailsServiceImpl(userAuthDao dao.UserAuthDao, userInfoDao dao.UserInfoDao, roleDao dao.RoleDao, redisService service.RedisService, store *sessions.CookieStore) *UserDetailsServiceImpl {
	return &UserDetailsServiceImpl{
		userAuthDao:  userAuthDao,
		userInfoDao:  userInfoDao,
		roleDao:      roleDao,
		redisService: redisService,
		store:        store,
	}
}

func (s *UserDetailsServiceImpl) LoadUserByUsername(r *http.Request, ctx context.Context, username string) (*dto.UserDetailDTO, error) {
	if username == "" {
		return nil, exception.NewBizError(enums.VALID_ERROR.Code, "用户名不能为空!")
	}

	userAuth, err := s.userAuthDao.SelectUserByUsername(ctx, username)
	if err != nil {
		return nil, exception.NewBizError(enums.SYSTEM_ERROR.Code, err.Error())
	}
	if userAuth == nil {
		return nil, exception.NewBizError(enums.USERNAME_NOT_EXIST.Code, "用户名不存在!")
	}

	return s.ConvertUserDetail(ctx, userAuth)
}

// LoadUserByID 根据用户ID加载用户信息
func (s *UserDetailsServiceImpl) LoadUserByID(ctx context.Context, userID int) (*dto.UserDetailDTO, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	if userID <= 0 {
		return nil, exception.NewBizError(enums.VALID_ERROR.Code, "Invalid user ID")
	}

	userAuth, err := s.userAuthDao.FindUserById(ctx, userID)
	if err != nil {
		return nil, exception.NewBizError(enums.SYSTEM_ERROR.Code, err.Error())
	}
	if userAuth == nil {
		return nil, exception.NewBizError(enums.USERNAME_NOT_EXIST.Code, "User does not exist")
	}

	userDetail, err := s.ConvertUserDetail(ctx, userAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user detail: %w", err)
	}

	return userDetail, nil
}

// ConvertUserDetail 将用户信息转换为 DTO
func (s *UserDetailsServiceImpl) ConvertUserDetail(ctx context.Context, user *model.UserAuth) (*dto.UserDetailDTO, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	userInfo, err := s.userInfoDao.GetUserInfoById(ctx, user.UserInfoID)
	if err != nil {
		return nil, exception.NewBizError(enums.SYSTEM_ERROR.Code, "获取用户信息失败")
	}

	roleList, err := s.roleDao.ListRolesByUserInfoId(ctx, userInfo.ID)
	if err != nil {
		return nil, exception.NewBizError(enums.SYSTEM_ERROR.Code, "获取角色列表失败")
	}

	// 获取并转换点赞数据为 []string 类型
	articleLikeSet, err := s.redisService.SMembers(ctx, constants.ArticleUserLike+strconv.Itoa(int(userInfo.ID)))
	if err != nil {
		return nil, exception.NewBizError(enums.SYSTEM_ERROR.Code, "获取文章点赞数据失败")
	}
	articleLikes := make([]string, len(articleLikeSet))
	copy(articleLikes, articleLikeSet)

	commentLikeSet, err := s.redisService.SMembers(ctx, constants.CommentUserLike+strconv.Itoa(int(userInfo.ID)))
	if err != nil {
		return nil, exception.NewBizError(enums.SYSTEM_ERROR.Code, "获取评论点赞数据失败")
	}
	commentLikes := make([]string, len(commentLikeSet))
	copy(commentLikes, commentLikeSet)

	talkLikeSet, err := s.redisService.SMembers(ctx, constants.TalkUserLike+strconv.Itoa(int(userInfo.ID)))
	if err != nil {
		return nil, exception.NewBizError(enums.SYSTEM_ERROR.Code, "获取讨论点赞数据失败")
	}
	talkLikes := make([]string, len(talkLikeSet))
	copy(talkLikes, talkLikeSet)

	// 构建 UserDetailDTO 并返回
	return &dto.UserDetailDTO{
		ID:             user.ID,
		LoginType:      user.LoginType,
		UserInfoID:     userInfo.ID,
		Username:       user.Username,
		Password:       user.Password,
		Email:          userInfo.Email,
		RoleList:       roleList,
		Nickname:       userInfo.Nickname,
		Avatar:         userInfo.Avatar,
		Intro:          userInfo.Intro,
		WebSite:        userInfo.WebSite,
		ArticleLikeSet: articleLikes,
		CommentLikeSet: commentLikes,
		TalkLikeSet:    talkLikes,
		IPAddress:      utils.GetIPAddressFromContext(ctx),
		IPSource:       utils.GetIPSource(utils.GetIPAddressFromContext(ctx)),
		//Browser:        utils.GetUserAgentFromContext(ctx).Browser(),
		OS:            utils.GetUserAgentFromContext(ctx).OS(),
		LastLoginTime: time.Now().In(time.FixedZone("Asia/Shanghai", 8*3600)),
	}, nil
}

func (s *UserDetailsServiceImpl) UpdateUserInfo(r *http.Request, userDetail *dto.UserDetailDTO) {
	userAuth := &model.UserAuth{
		ID:            userDetail.ID,
		IPAddress:     userDetail.IPAddress,
		IPSource:      userDetail.IPSource,
		LastLoginTime: userDetail.LastLoginTime,
	}
	s.userAuthDao.UpdateById(r.Context(), userAuth)
}
