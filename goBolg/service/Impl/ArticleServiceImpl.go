package Impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/enums"
	"goBolg/model"
	"goBolg/service"
	"goBolg/strategy/contxt"
	"goBolg/utils"
	"goBolg/vo"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
)

type ArticleServiceImpl struct {
	articleDao      dao.ArticleDao
	categoryDao     dao.CategoryDao
	tagDao          dao.TagDao
	tagService      service.TagService
	articleTagDao   dao.ArticleTagDao
	redisService    service.RedisService
	searchStrategy  *contxt.SearchStrategyContext
	blogInfoService service.BlogInfoService
	articleSet      sync.Map // 使用 sync.Map 代替 Set
	db              *gorm.DB // 添加对 gorm.DB 的引用，用于事务处理
}

func NewArticleServiceImpl(articleDao dao.ArticleDao, articleTagDao dao.ArticleTagDao, categoryDao dao.CategoryDao, tagDao dao.TagDao, tagService service.TagService, redisService service.RedisService, blogInfoService service.BlogInfoService, searchStrategy *contxt.SearchStrategyContext, db *gorm.DB) *ArticleServiceImpl {
	return &ArticleServiceImpl{
		articleDao:      articleDao,
		categoryDao:     categoryDao,
		tagDao:          tagDao,
		tagService:      tagService,
		articleTagDao:   articleTagDao,
		redisService:    redisService,
		blogInfoService: blogInfoService,
		searchStrategy:  searchStrategy,
		db:              db,
	}
}

// ListArchives 查询归档列表
func (s *ArticleServiceImpl) ListArchives(current int) ([]dto.ArchiveDTO, int, error) {
	log.Println("Service: ListArchives called")

	// 设置分页参数
	limit := 10 // 设置每页显示的条目数
	offset := (current - 1) * limit

	archives, err := s.articleDao.ListArchives(offset, limit)
	if err != nil {
		log.Printf("Error in rabbitService layer: %v", err)
		return nil, 0, err
	}

	totalCount, err := s.articleDao.CountArchives()
	if err != nil {
		log.Printf("Error in counting archives: %v", err)
		return nil, 0, err
	}

	log.Printf("Service: ListArchives returned %d records", len(archives))
	return archives, totalCount, nil
}

func (s *ArticleServiceImpl) ListArticles(ctx context.Context) ([]dto.ArticleHomeDTO, error) {
	offset := utils.GetLimitCurrent(ctx)
	size := utils.GetSize(ctx)
	return s.articleDao.ListArticles(offset, size)
}

func (s *ArticleServiceImpl) CountArticles(ctx context.Context) (int64, error) {
	return s.articleDao.CountArticles(ctx)
}

func (s *ArticleServiceImpl) ListArticleBacks(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.ArticleBackDTO], error) {
	// 获取分页参数
	limit := utils.GetSize(ctx)
	offset := utils.GetLimitCurrent(ctx)

	// 查询文章总量
	count, err := s.articleDao.CountArticleBacks(ctx, condition)
	if err != nil || count == 0 {
		return vo.NewPageResult([]dto.ArticleBackDTO{}, 0), err
	}

	// 查询后台文章并预加载标签信息
	articleBackDTOList, err := s.articleDao.ListArticleBacks(ctx, offset, limit, condition)
	if err != nil {
		return vo.NewPageResult([]dto.ArticleBackDTO{}, 0), err
	}

	// 查询文章点赞量和浏览量
	viewsCountMap, err := s.redisService.ZAllScore(ctx, constants.ArticleViewsCount)
	if err != nil {
		return vo.NewPageResult([]dto.ArticleBackDTO{}, 0), err
	}

	likeCountMap, err := s.redisService.HGetAll(ctx, constants.ArticleLikeCount)
	if err != nil {
		return vo.NewPageResult([]dto.ArticleBackDTO{}, 0), err
	}

	// 封装点赞量和浏览量
	for i := range articleBackDTOList {
		item := &articleBackDTOList[i]

		// 将 item.ID 转换为字符串进行查找
		idStr := strconv.Itoa(item.ID)

		// 处理 viewsCount
		if viewsCount, exists := viewsCountMap[idStr]; exists {
			item.ViewsCount = int(viewsCount)
		}

		// 处理 likeCount
		if likeCountStr, exists := likeCountMap[idStr]; exists {
			if likeCountInt, err := strconv.Atoi(likeCountStr); err == nil {
				item.LikeCount = likeCountInt
			}
		}

		// 查询并填充标签信息
		tags, err := s.articleDao.GetTagsByArticleID(ctx, item.ID)
		if err == nil {
			item.TagDTOList = tags
		}
	}

	return vo.NewPageResult(articleBackDTOList, count), nil
}

// GetArticleById 根据文章 ID 获取文章详情
func (s *ArticleServiceImpl) GetArticleById(ctx context.Context, articleId int) (*dto.ArticleDTO, error) {
	var wg sync.WaitGroup

	// 查询推荐文章
	var recommendArticleList []dto.ArticleRecommendDTO
	var recommendArticleErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		recommendArticleList, recommendArticleErr = s.articleDao.ListRecommendArticles(ctx, articleId)
	}()

	// 查询最新文章
	var newestArticleList []dto.ArticleRecommendDTO
	var newestArticleErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		articleList, err := s.articleDao.ListArticles(1, 5) // 获取最新文章，假设使用默认的当前页 1 和页大小 5
		if err != nil {
			newestArticleErr = err
			return
		}
		newestArticleList = utils.BeanCopyList(articleList, &dto.ArticleRecommendDTO{}).([]dto.ArticleRecommendDTO)
	}()

	// 查询指定文章
	article, err := s.articleDao.GetArticleById(ctx, articleId)
	if err != nil {
		return nil, err
	}

	if article.ID == 0 {
		return nil, errors.New("文章不存在")
	}
	// 获取用户信息
	user, ok := utils.GetLoginUser(ctx)
	if ok && user != nil {
		userInfoID := user.UserInfoID
		articleLikeKey := constants.ArticleUserLike + strconv.Itoa(userInfoID)
		likedArticles, err := s.redisService.SMembers(ctx, articleLikeKey)
		if err != nil {
			return nil, err
		}
		liked := false
		for _, likedArticle := range likedArticles {
			if likedArticle == strconv.Itoa(articleId) {
				liked = true
				break
			}
		}
		if liked {
			article.LikeCount++ // 增加 1 表示用户点赞
		}
	}

	// 更新文章浏览量
	s.UpdateArticleViewsCount(ctx, articleId)
	// 获取浏览量
	viewsCountStr, err := s.redisService.ZScore(ctx, constants.ArticleViewsCount, strconv.Itoa(articleId))
	if err != nil {
		log.Printf("Error retrieving views count: %v", err)
		return nil, err
	}

	// 处理可能的 nil 值
	viewsCount := 0
	if viewsCountStr != 0 {
		viewsCount = int(viewsCountStr)
	}
	article.ViewsCount = viewsCount
	// 查询上一篇和下一篇文章
	lastArticle, err := s.articleDao.GetPreviousArticle(ctx, articleId)
	if err != nil {
		return nil, err
	}
	nextArticle, err := s.articleDao.GetNextArticle(ctx, articleId)
	if err != nil {
		return nil, err
	}

	article.LastArticle = &dto.ArticlePaginationDTO{}
	article.NextArticle = &dto.ArticlePaginationDTO{}
	utils.BeanCopy(&lastArticle, article.LastArticle)
	utils.BeanCopy(&nextArticle, article.NextArticle)

	// 获取点赞次数
	likeCountStr, err := s.redisService.HGet(ctx, constants.ArticleLikeCount, strconv.Itoa(articleId))
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		article.LikeCount = 0
		_, err = s.redisService.HSet(ctx, constants.ArticleLikeCount, strconv.Itoa(articleId), "0")
		if err != nil {
			return nil, err
		}
	} else {
		likeCount, err := strconv.Atoi(likeCountStr)
		if err != nil {
			return nil, err
		}
		article.LikeCount = likeCount
	}

	// 等待所有并发操作完成
	wg.Wait()

	// 封装推荐文章和最新文章
	if recommendArticleErr != nil {
		log.Printf("Error fetching recommend articles: %v", recommendArticleErr)
	} else {
		article.RecommendArticleList = recommendArticleList
	}

	if newestArticleErr != nil {
		log.Printf("Error fetching newest articles: %v", newestArticleErr)
	} else {
		article.NewestArticleList = newestArticleList
	}

	return &article, nil
}

// UpdateArticleViewsCount 更新文章浏览量
func (s *ArticleServiceImpl) UpdateArticleViewsCount(ctx context.Context, articleId int) {
	// 检查 articleId 是否已被访问
	if _, loaded := s.articleSet.LoadOrStore(articleId, true); !loaded {
		// 浏览量 +1
		_, err := s.redisService.ZIncrBy(ctx, constants.ArticleViewsCount, strconv.Itoa(articleId), 1.0)
		if err != nil {
			log.Printf("Error incrementing views count for article %d: %v", articleId, err)
		}
	}
}

func (s *ArticleServiceImpl) ListArticlesByCondition(ctx context.Context, condition vo.ConditionVO) (*dto.ArticlePreviewListDTO, error) {
	// 获取分页参数
	limit := utils.GetSize(ctx)
	offset := utils.GetLimitCurrent(ctx)

	// 查询文章
	articles, err := s.articleDao.ListArticlesByCondition(ctx, offset, limit, condition)
	if err != nil {
		return nil, err
	}

	// 根据条件获取分类名或标签名
	var name string
	if condition.CategoryID != nil {
		category, err := s.categoryDao.GetCategoryNameByID(ctx, *condition.CategoryID)
		if err != nil {
			return nil, err
		}
		name = category
	} else if condition.TagID != nil {
		tag, err := s.tagDao.GetTagNameByID(ctx, *condition.TagID)
		if err != nil {
			return nil, err
		}
		name = tag
	} else {
		name = "所有文章"
	}

	return &dto.ArticlePreviewListDTO{
		ArticlePreviewDTOList: articles,
		Name:                  name,
	}, nil
}

// GetArticleBackById 根据ID获取后台文章
func (s *ArticleServiceImpl) GetArticleBackById(ctx context.Context, articleId int) (*vo.ArticleVO, error) {
	// 查询文章信息
	article, err := s.articleDao.GetArticleById(ctx, articleId)
	if err != nil {
		return nil, err
	}
	if article.ID == 0 {
		return nil, errors.New("文章不存在")
	}

	// 查询文章分类
	category, err := s.categoryDao.GetCategoryByID(ctx, article.CategoryId)
	if err != nil {
		return nil, err
	}
	var categoryName *string
	if category != nil {
		categoryName = &category.Name
	}

	// 查询文章标签
	tagNameList, err := s.tagDao.ListTagNameByArticleId(ctx, articleId)
	if err != nil {
		return nil, err
	}

	// 封装数据
	articleVO := &vo.ArticleVO{
		ID:             &article.ID,
		ArticleTitle:   article.ArticleTitle,
		ArticleContent: article.ArticleContent,
		ArticleCover:   &article.ArticleCover,
		CategoryName:   categoryName,
		TagNameList:    tagNameList,
		Type:           &article.Type,
		OriginalURL:    &article.OriginalUrl,
	}

	return articleVO, nil
}
func (s *ArticleServiceImpl) SaveArticleLike(ctx context.Context, articleId int) (int, error) {
	user, ok := utils.GetLoginUser(ctx)
	if !ok || user == nil {
		log.Printf("User not found in context. Context: %+v", ctx)
		return 0, fmt.Errorf("user not found in context")
	}

	userInfoID := user.UserInfoID
	articleLikeKey := constants.ArticleUserLike + strconv.Itoa(userInfoID)
	articleIdStr := strconv.Itoa(articleId)

	// 获取当前用户的点赞文章列表
	liked, err := s.redisService.SMembers(ctx, articleLikeKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get liked articles: %v", err)
	}

	likedArticles := make(map[int]struct{})
	for _, id := range liked {
		if articleIdInt, err := strconv.Atoi(id); err == nil {
			likedArticles[articleIdInt] = struct{}{}
		}
	}

	// 判断用户是否已经点赞过该文章
	if _, exists := likedArticles[articleId]; exists {
		// 取消点赞
		removed, err := s.redisService.SRemove(ctx, articleLikeKey, articleIdStr)
		if err != nil {
			return 0, fmt.Errorf("failed to remove article like: %v", err)
		}
		if removed == 0 {
			log.Printf("Article %d not found in user's likes, cannot remove", articleId)
		} else {
			log.Printf("Article %d successfully removed from user's likes", articleId)
		}

		// 减少文章点赞数量
		_, err = s.redisService.HDecr(ctx, constants.ArticleLikeCount, articleIdStr, 1)
		if err != nil {
			return 0, fmt.Errorf("failed to decrement like count: %v", err)
		}
		log.Printf("Article %d unliked by user %d", articleId, userInfoID)
	} else {
		// 添加点赞
		added, err := s.redisService.SAdd(ctx, articleLikeKey, articleIdStr)
		if err != nil {
			return 0, fmt.Errorf("failed to add article like: %v", err)
		}
		if added == 0 {
			log.Printf("Article %d already liked by user %d", articleId, userInfoID)
		} else {
			log.Printf("Article %d successfully added to user's likes", articleId)
		}

		// 增加文章点赞数量
		_, err = s.redisService.HIncr(ctx, constants.ArticleLikeCount, articleIdStr, 1)
		if err != nil {
			return 0, fmt.Errorf("failed to increment like count: %v", err)
		}
		log.Printf("Article %d liked by user %d", articleId, userInfoID)
	}

	// 检查最终的点赞数量
	likeCount, err := s.redisService.HGet(ctx, constants.ArticleLikeCount, articleIdStr)
	if err != nil {
		log.Printf("Failed to get final like count for article %d: %v", articleId, err)
		return 0, fmt.Errorf("failed to get final like count: %v", err)
	}

	finalLikeCount, err := strconv.Atoi(likeCount)
	if err != nil {
		log.Printf("Failed to convert like count for article %d: %v", articleId, err)
		return 0, fmt.Errorf("failed to convert like count: %v", err)
	}

	log.Printf("Final like count for article %d: %d", articleId, finalLikeCount)

	return finalLikeCount, nil
}

// SaveOrUpdateArticle 保存或更新文章
func (s *ArticleServiceImpl) SaveOrUpdateArticle(ctx context.Context, articleVO vo.ArticleVO) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var webConfig vo.WebsiteConfigVO
	var category *model.Category
	var err error

	// 并发获取博客配置信息
	wg.Add(1)
	go func() {
		defer wg.Done()
		config, e := s.blogInfoService.GetWebsiteConfig(ctx)
		mu.Lock()
		defer mu.Unlock()
		if e != nil {
			err = e
			return
		}
		webConfig = config
	}()

	// 保存文章分类
	category, err = s.saveArticleCategory(ctx, articleVO)
	if err != nil {

		return err
	}

	// 等待并发操作完成
	wg.Wait()
	if err != nil {
		return err
	}

	// 获取当前登录用户
	user, ok := utils.GetLoginUser(ctx)
	if !ok {
		return errors.New("failed to get login user")
	}

	article := model.Article{}
	utils.BeanCopy(&articleVO, &article)

	if category != nil {
		article.CategoryID = category.ID
	}

	// 检查文章封面是否为空
	if articleVO.ArticleCover != nil && *articleVO.ArticleCover != "" {
		article.ArticleCover = *articleVO.ArticleCover
	} else {
		article.ArticleCover = webConfig.ArticleCover
	}

	article.UserID = user.UserInfoID

	// 解引用 ArticleVO 中的 ID 并判断是否为 0
	if articleVO.ID == nil || *articleVO.ID == 0 {
		err = s.articleDao.SaveArticle(ctx, &article)
	} else {
		article.ID = *articleVO.ID // 将解引用后的 ID 赋值给 article.ID
		err = s.articleDao.UpdateArticle(ctx, &article)
	}

	if err != nil {
		return err
	}

	// 保存文章标签
	err = s.saveArticleTag(ctx, articleVO, article.ID)
	if err != nil {
		return err
	}

	return nil
}

// SaveArticleTag 保存文章标签
func (s *ArticleServiceImpl) saveArticleTag(ctx context.Context, articleVO vo.ArticleVO, articleId int) error {
	tagNameList := articleVO.TagNameList
	if len(tagNameList) > 0 {
		// 查询数据库中已存在的标签
		existTagList, err := s.tagDao.ListTagsByNames(ctx, tagNameList)
		if err != nil {
			return err
		}

		existTagNameMap := make(map[string]int)
		for _, tag := range existTagList {
			existTagNameMap[tag.TagName] = tag.ID
		}

		// 准备保存的新文章-标签关联关系
		var newArticleTagList []model.ArticleTag

		for _, tagID := range existTagNameMap {
			exists, err := s.articleTagDao.ExistsByArticleAndTag(ctx, articleId, tagID)
			if err != nil {
				return err
			}

			if exists {
				err := s.articleTagDao.UpdateTimestamp(ctx, articleId, tagID, time.Now())
				if err != nil {
					return err
				}
			} else {

				newArticleTagList = append(newArticleTagList, model.ArticleTag{
					ArticleID: articleId,
					TagID:     tagID,
				})
			}
		}

		// 保存新的文章-标签关联关系
		if len(newArticleTagList) > 0 {
			log.Printf("Saving new article-tag associations: %v", newArticleTagList)
			err := s.articleTagDao.SaveBatch(newArticleTagList)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ArticleServiceImpl) saveArticleCategory(ctx context.Context, articleVO vo.ArticleVO) (*model.Category, error) {
	// 检查 articleVO.CategoryName 是否为 nil
	if articleVO.CategoryName == nil {
		return nil, errors.New("category name is required")
	}

	// 查询分类是否存在
	categoryDTO, err := s.categoryDao.GetCategoryByName(ctx, *articleVO.CategoryName)
	if err != nil {
		return nil, err
	}

	// 如果分类不存在且文章状态不为草稿，则创建新分类
	if categoryDTO == nil && (articleVO.Status == nil || *articleVO.Status != enums.DRAFT.Status) {
		newCategory := &model.Category{
			Name: *articleVO.CategoryName,
		}
		if err := s.categoryDao.InsertCategory(ctx, newCategory); err != nil {
			return nil, err
		}
		return newCategory, nil
	}

	// 返回已存在的分类
	if categoryDTO != nil {
		return &model.Category{
			ID:   categoryDTO.ID,
			Name: categoryDTO.Name,
		}, nil
	}

	return nil, errors.New("unable to create or find category")
}

// UpdateArticleTop 修改文章置顶状态
func (s *ArticleServiceImpl) UpdateArticleTop(ctx context.Context, articleTopVO vo.ArticleTopVO) error {
	if articleTopVO.IsTop == nil {
		return fmt.Errorf("IsTop is nil")
	}

	article := &model.Article{
		ID:    *articleTopVO.ID,
		IsTop: articleTopVO.IsTop,
	}

	log.Printf("Updating article with ID: %d, IsTop: %d", article.ID, *article.IsTop) // 打印更新的信息

	return s.articleDao.UpdateById(ctx, article)
}

// UpdateArticleDelete 修改文章逻辑删除状态
func (s *ArticleServiceImpl) UpdateArticleDelete(ctx context.Context, deleteVO vo.DeleteVO) error {
	var articles []model.Article
	falseInt := int(constants.False)
	pFalse := &falseInt
	for _, id := range deleteVO.IDList {
		articles = append(articles, model.Article{
			ID:       id,
			IsTop:    pFalse,
			IsDelete: deleteVO.IsDelete,
		})
	}

	// 批量更新文章状态
	if err := s.articleDao.UpdateBatchById(ctx, articles); err != nil {
		return fmt.Errorf("failed to update articles: %v", err)
	}
	return nil
}

// DeleteArticles 删除文章和相关的文章标签
func (s *ArticleServiceImpl) DeleteArticles(ctx context.Context, articleIdList []int) error {
	// 启动事务
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("启动事务失败: %v", tx.Error)
	}
	log.Println("事务启动成功")
	// 确保事务要么提交要么回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("事务回滚: %v", r)
		}
	}()

	// 删除文章标签关联
	if err := s.articleTagDao.DeleteByArticleIds(ctx, articleIdList); err != nil {
		tx.Rollback()
		return fmt.Errorf("删除文章标签失败: %v", err)
	}

	// 删除文章
	if err := s.articleDao.DeleteArticlesByIds(ctx, articleIdList); err != nil {
		tx.Rollback()
		return fmt.Errorf("删除文章失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("事务提交失败: %v", err)
	}

	return nil
}

func (s *ArticleServiceImpl) ListArticlesBySearch(ctx context.Context, condition vo.ConditionVO) ([]dto.ArticleSearchDTO, error) {
	var keywords string
	if condition.Keywords != nil {
		keywords = *condition.Keywords
	}
	return s.searchStrategy.ExecuteSearchStrategy(ctx, keywords)
}
