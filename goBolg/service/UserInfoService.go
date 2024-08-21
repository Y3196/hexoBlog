package service

import (
	"context"
	"goBolg/vo"
	"mime/multipart"
)

// UserInfoService 接口定义
type UserInfoService interface {
	UpdateUserInfo(ctx context.Context, userInfoVO *vo.UserInfoVO) error

	UpdateUserAvatar(ctx context.Context, file *multipart.FileHeader) (string, error)

	SaveUserEmail(ctx context.Context, emailVO vo.EmailVO) error
}
