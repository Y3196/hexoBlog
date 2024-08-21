package service

import (
	"context"
	"goBolg/dto"
)

// UniqueViewService defines the rabbitService for managing unique views.
type UniqueViewService interface {
	// ListUniqueViews retrieves statistics for unique views over the last 7 days.
	ListUniqueViews(ctx context.Context) ([]dto.UniqueViewDTO, error)
}
