package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

type MenuService interface {
	ListMenus(ctx context.Context, conditionVO vo.ConditionVO) ([]dto.MenuDTO, error)

	SaveOrUpdateMenu(ctx context.Context, menuVO vo.MenuVO) error

	DeleteMenu(ctx context.Context, menuId uint) error

	ListMenuOptions(ctx context.Context) ([]dto.LabelOptionDTO, error)

	ListUserMenus(ctx context.Context, userID uint) ([]dto.UserMenuDTO, error)
}
