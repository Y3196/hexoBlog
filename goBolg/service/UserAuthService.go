package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

// UserAuthService 是用户认证服务接口
type UserAuthService interface {
	SendCode(username string) error

	Register(ctx context.Context, user vo.UserVO) error

	UpdatePassword(ctx context.Context, user *vo.UserVO) error

	ListUserAreas(ctx context.Context, conditionVO vo.ConditionVO) ([]dto.UserAreaDTO, error)
}
