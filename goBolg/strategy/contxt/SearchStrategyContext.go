package contxt

import (
	"context"
	"fmt"
	"goBolg/dto"
	"goBolg/strategy"
)

// SearchStrategyContext 搜索策略上下文
type SearchStrategyContext struct {
	searchMode        string
	searchStrategyMap map[string]strategy.SearchStrategy
}

// NewSearchStrategyContext 创建SearchStrategyContext实例
func NewSearchStrategyContext(searchMode string, searchStrategyMap map[string]strategy.SearchStrategy) *SearchStrategyContext {
	return &SearchStrategyContext{
		searchMode:        searchMode,
		searchStrategyMap: searchStrategyMap,
	}
}

// ExecuteSearchStrategy 执行搜索策略
func (c *SearchStrategyContext) ExecuteSearchStrategy(ctx context.Context, keywords string) ([]dto.ArticleSearchDTO, error) {
	strategy, exists := c.searchStrategyMap[c.searchMode]
	if !exists {
		return nil, fmt.Errorf("search strategy not found for mode: %s", c.searchMode)
	}
	return strategy.SearchArticle(ctx, keywords)
}
