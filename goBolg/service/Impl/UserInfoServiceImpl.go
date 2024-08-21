package Impl

import (
	"context"
	"errors"
	"fmt"
	constants "goBolg/constant"
	"goBolg/dao"
	"goBolg/enums"
	"goBolg/model"
	"goBolg/service"
	"goBolg/strategy/contxt"
	"goBolg/utils"
	"goBolg/vo"
	"log"
	"mime/multipart"
)

// userInfoServiceImpl 实现了 UserInfoService 接口
type userInfoServiceImpl struct {
	userInfoDao           dao.UserInfoDao
	uploadStrategyContext *contxt.UploadStrategyContext
	redisService          service.RedisService
}

// NewUserInfoService 创建一个新的 UserInfoService 实例
func NewUserInfoService(userInfoDao dao.UserInfoDao, uploadStrategyContext *contxt.UploadStrategyContext, redisService service.RedisService) service.UserInfoService {
	return &userInfoServiceImpl{
		userInfoDao:           userInfoDao,
		uploadStrategyContext: uploadStrategyContext,
		redisService:          redisService,
	}
}

// UpdateUserInfo 更新用户信息
func (s *userInfoServiceImpl) UpdateUserInfo(ctx context.Context, userInfoVO *vo.UserInfoVO) error {
	// 从上下文获取当前登录的用户信息
	user, ok := utils.GetLoginUser(ctx)
	if !ok {
		return errors.New("failed to get login user from context")
	}

	// 封装用户信息
	userInfo := &model.UserInfo{
		ID:       user.UserInfoID, // 使用从上下文获取的用户ID
		Nickname: userInfoVO.Nickname,
		Intro:    userInfoVO.Intro,
		WebSite:  userInfoVO.WebSite,
	}

	// 更新用户信息
	return s.userInfoDao.UpdateUserInfo(ctx, userInfo)
}

func (s *userInfoServiceImpl) UpdateUserAvatar(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	// 获取用户信息
	user, ok := utils.GetLoginUser(ctx)
	if !ok {
		return "", fmt.Errorf("failed to retrieve user from context")
	}

	// 上传文件
	filePathEnum := enums.Avatar
	avatarURL, err := s.uploadStrategyContext.ExecuteUploadStrategy(fileHeader, filePathEnum)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// 更新用户信息
	userInfo := model.UserInfo{ID: user.ID, Avatar: avatarURL}
	if err := s.userInfoDao.UpdateUserInfo(ctx, &userInfo); err != nil {
		return "", fmt.Errorf("failed to update user info: %w", err)
	}

	return avatarURL, nil
}

// SaveUserEmail 更新用户邮箱
func (service *userInfoServiceImpl) SaveUserEmail(ctx context.Context, emailVO vo.EmailVO) error {
	// 验证验证码
	code, err := service.redisService.Get(ctx, constants.UserCodeKey+emailVO.Email)
	if err != nil {
		log.Printf("Failed to get code from Redis: %v", err)
		return errors.New("验证码错误")
	}
	if code != emailVO.Code {
		return errors.New("验证码错误")
	}

	// 获取当前登录用户的 ID
	userDetailDTO, ok := utils.GetLoginUser(ctx)
	if !ok {
		log.Println("Failed to get user from context")
		return errors.New("获取用户信息失败")
	}

	userInfoId := userDetailDTO.UserInfoID
	log.Printf("Updating email for user ID: %d", userInfoId)

	// 更新用户邮箱
	userInfo := &model.UserInfo{
		ID:    userInfoId,
		Email: emailVO.Email,
	}
	err = service.userInfoDao.UpdateUserInfo(ctx, userInfo)
	if err != nil {
		log.Printf("Failed to update user info: %v", err)
		return err
	}

	return nil
}
