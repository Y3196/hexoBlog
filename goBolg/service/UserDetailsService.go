package service

import (
	"context"
	"goBolg/dto"
	"net/http"
)

type UserDetailsService interface {
	LoadUserByUsername(r *http.Request, ctx context.Context, username string) (*dto.UserDetailDTO, error)

	LoadUserByID(ctx context.Context, userID int) (*dto.UserDetailDTO, error)

	UpdateUserInfo(r *http.Request, userDetail *dto.UserDetailDTO)
}
