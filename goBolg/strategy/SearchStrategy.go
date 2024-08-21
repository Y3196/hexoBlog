package strategy

import (
	"context"
	"goBolg/dto"
)

// SearchStrategy 搜索策略接口
type SearchStrategy interface {
	SearchArticle(ctx context.Context, keywords string) ([]dto.ArticleSearchDTO, error)
}
