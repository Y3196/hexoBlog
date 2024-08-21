package service

import (
	"context"
	"goBolg/dto"
	"goBolg/vo"
)

// ArticleService 文章服务接口
type ArticleService interface {
	ListArchives(current int) ([]dto.ArchiveDTO, int, error)

	ListArticles(ctx context.Context) ([]dto.ArticleHomeDTO, error)

	CountArticles(ctx context.Context) (int64, error)

	// 查询后台文章
	ListArticleBacks(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.ArticleBackDTO], error)

	//根据id查看文章
	GetArticleById(ctx context.Context, articleId int) (*dto.ArticleDTO, error)

	// 根据条件查询文章列表
	ListArticlesByCondition(ctx context.Context, condition vo.ConditionVO) (*dto.ArticlePreviewListDTO, error)

	// 根据id查看后台文章
	GetArticleBackById(ctx context.Context, articleId int) (*vo.ArticleVO, error)

	// 点赞文章
	SaveArticleLike(ctx context.Context, articleId int) (int, error)

	// 添加或修改文章
	SaveOrUpdateArticle(ctx context.Context, articleVO vo.ArticleVO) error

	// 修改文章置顶
	UpdateArticleTop(ctx context.Context, articleTopVO vo.ArticleTopVO) error

	// 删除或恢复文章
	UpdateArticleDelete(ctx context.Context, deleteVO vo.DeleteVO) error

	// 物理删除文章
	DeleteArticles(ctx context.Context, articleIdList []int) error

	// 搜索文章
	ListArticlesBySearch(ctx context.Context, condition vo.ConditionVO) ([]dto.ArticleSearchDTO, error)
	/*


		// 导出文章
		ExportArticles(ctx contxt.Context, articleIdList []int) ([]string, error)*/
}
