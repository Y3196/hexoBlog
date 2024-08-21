package utils

import (
	"context"
	"goBolg/constant"
	"goBolg/dto"
	"log"
)

func GetLoginUser(ctx context.Context) (*dto.UserDetailDTO, bool) {
	user, ok := ctx.Value(constants.UserContextKey).(*dto.UserDetailDTO)
	log.Printf("Getting user from context: %+v, found: %v", user, ok)
	return user, ok
}
