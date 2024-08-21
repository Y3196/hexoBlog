package strategyImpl

import (
	"context"
	"fmt"
	"goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/enums"
	"strings"
)

// MySqlSearchStrategyImpl MySQL搜索策略实现
type MySqlSearchStrategyImpl struct {
	ArticleDao dao.ArticleDao
}

// SearchArticle 搜索文章
func (s *MySqlSearchStrategyImpl) SearchArticle(ctx context.Context, keywords string) ([]dto.ArticleSearchDTO, error) {
	if keywords == "" {
		return []dto.ArticleSearchDTO{}, nil
	}

	var articles []dto.ArticleSearchDTO

	// 查询文章
	err := s.ArticleDao.GetDb().WithContext(ctx).Table("tb_article").
		Select("id, article_title, article_content").
		Where("is_delete = ? AND status = ? AND (article_title LIKE ? OR article_content LIKE ?)",
			false, enums.PUBLIC.Status, fmt.Sprintf("%%%s%%", keywords), fmt.Sprintf("%%%s%%", keywords)).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	// 高亮处理
	for i, article := range articles {
		// 文章内容高亮
		articleContent := article.ArticleContent
		index := strings.Index(articleContent, keywords)
		if index != -1 {
			preIndex := max(0, index-25)
			postIndex := min(len(articleContent), index+len(keywords)+175)
			highlightedContent := articleContent[preIndex:postIndex]
			articles[i].ArticleContent = strings.ReplaceAll(highlightedContent, keywords, fmt.Sprintf("%s%s%s", constants.PreTag, keywords, constants.PostTag))
		}
		// 文章标题高亮
		articles[i].ArticleTitle = strings.ReplaceAll(article.ArticleTitle, keywords, fmt.Sprintf("%s%s%s", constants.PreTag, keywords, constants.PostTag))
	}

	return articles, nil
}

// max 返回两个整数中的最大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min 返回两个整数中的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
