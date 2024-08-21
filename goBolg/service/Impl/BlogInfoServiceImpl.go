package Impl

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type blogInfoServiceImpl struct {
	db               *gorm.DB
	userInfoDao      dao.UserInfoDao
	messageDao       dao.MessageDao
	articleDao       dao.ArticleDao
	categoryDao      dao.CategoryDao
	tagDao           dao.TagDao
	pageService      service.PageService
	redisService     service.RedisService
	websiteConfigDao dao.WebsiteConfigDao
	uniqueViewDao    dao.UniqueViewDao
}

func NewBlogInfoService(userInfoDao dao.UserInfoDao, messageDao dao.MessageDao, uniqueViewDao dao.UniqueViewDao, articleDao dao.ArticleDao, categoryDao dao.CategoryDao, tagDao dao.TagDao, redisService *RedisServiceImpl, websiteConfigDao dao.WebsiteConfigDao, pageService service.PageService, db *gorm.DB) service.BlogInfoService {
	return &blogInfoServiceImpl{
		db:               db,
		userInfoDao:      userInfoDao,
		messageDao:       messageDao,
		uniqueViewDao:    uniqueViewDao,
		articleDao:       articleDao,
		categoryDao:      categoryDao,
		tagDao:           tagDao,
		redisService:     redisService,
		websiteConfigDao: websiteConfigDao,
		pageService:      pageService,
	}
}

func (s *blogInfoServiceImpl) GetBlogHomeInfo(ctx context.Context) (dto.BlogHomeInfoDTO, error) {
	fmt.Println("Entering GetBlogHomeInfo method")

	articleCount, err := s.articleDao.CountArticles(ctx)
	if err != nil {
		fmt.Printf("Error counting articles: %v\n", err)
		return dto.BlogHomeInfoDTO{}, err
	}
	fmt.Printf("Article count: %d\n", articleCount)

	categoryCount, err := s.categoryDao.CountCategories(ctx)
	if err != nil {
		fmt.Printf("Error counting categories: %v\n", err)
		return dto.BlogHomeInfoDTO{}, err
	}
	fmt.Printf("Category count: %d\n", categoryCount)

	tagCount, err := s.tagDao.CountTags(ctx)
	if err != nil {
		fmt.Printf("Error counting tags: %v\n", err)
		return dto.BlogHomeInfoDTO{}, err
	}
	fmt.Printf("Tag count: %d\n", tagCount)

	viewsCount, err := s.redisService.Get(ctx, constants.BlogViewsCount)
	if err != nil {
		if err == redis.Nil {
			viewsCount = "0"
			fmt.Println("Views count not found in Redis, defaulting to 0")
		} else {
			fmt.Printf("Error retrieving views count from Redis: %v\n", err)
			return dto.BlogHomeInfoDTO{}, err
		}
	}
	fmt.Printf("Views count from Redis: %s\n", viewsCount)

	if viewsCount == "" {
		viewsCount = "0"
		fmt.Println("Views count is empty, defaulting to 0")
	}

	websiteConfig, err := s.GetWebsiteConfig(ctx)
	if err != nil {
		fmt.Printf("Error retrieving website configuration: %v\n", err)
		return dto.BlogHomeInfoDTO{}, err
	}
	fmt.Printf("Website configuration: %+v\n", websiteConfig)

	pageList, err := s.pageService.ListPages(ctx)
	if err != nil {
		fmt.Printf("Error listing pages: %v\n", err)
		return dto.BlogHomeInfoDTO{}, err
	}
	fmt.Printf("Page list: %+v\n", pageList)

	fmt.Println("Exiting GetBlogHomeInfo method")

	return *dto.NewBlogHomeInfoDTO(
		int(articleCount),
		int(categoryCount),
		int(tagCount),
		viewsCount,
		websiteConfig,
		pageList,
	), nil
}

func (s *blogInfoServiceImpl) GetWebsiteConfig(ctx context.Context) (vo.WebsiteConfigVO, error) {
	log.Println("获取网站配置")

	var configVO vo.WebsiteConfigVO
	var configStr string

	// 从 Redis 中获取 JSON 字符串
	configStr, err := s.redisService.Get(ctx, constants.WebsiteConfig)
	if err != nil {
		if err == redis.Nil {
			// Redis 中没有找到，需从数据库获取
			config, err := s.websiteConfigDao.SelectByID(ctx, s.db, constants.DefaultConfigID)
			if err != nil {
				return configVO, err
			}
			configStr = config.Config

			// 将数据存入 Redis 以备后用
			err = s.redisService.Set(ctx, constants.WebsiteConfig, configStr, time.Hour)
			if err != nil {
				log.Printf("存储数据到 Redis 失败: %v", err)
				return configVO, err
			}
		} else {
			return configVO, err
		}
	}

	// 尝试解码 JSON 字符串
	decodedStr, err := utils.DecodeDoubleEscapedJSON(configStr)
	if err != nil {
		log.Printf("解码 JSON 失败: %v", err)
		log.Printf("导致错误的 JSON 字符串: %s", configStr)
		return configVO, err
	}

	// 将解码后的 JSON 字符串反序列化到 configVO
	err = json.Unmarshal([]byte(decodedStr), &configVO)
	if err != nil {
		log.Printf("将 JSON 反序列化到 WebsiteConfigVO 失败: %v", err)
		log.Printf("导致错误的解码后 JSON 字符串: %s", decodedStr)
		return configVO, err
	}

	log.Printf("成功反序列化 WebsiteConfigVO: %+v", configVO)
	return configVO, nil
}

func (s *blogInfoServiceImpl) GetBlogBackInfo(ctx context.Context) (dto.BlogBackInfoDTO, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	// 获取 views count
	viewsCountStr, err := s.redisService.Get(ctx, "BLOG_VIEWS_COUNT")
	if err != nil && err != redis.Nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error retrieving views count: %w", err)
	}
	viewsCount, err := strconv.Atoi(utils.DefaultIfEmpty(viewsCountStr, "0"))
	if err != nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error converting views count: %w", err)
	}

	// 获取 message count
	messageCount, err := s.messageDao.CountMessages(ctx)
	if err != nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error counting messages: %w", err)
	}

	// 获取 user count
	userCount, err := s.userInfoDao.CountUsers(ctx)
	if err != nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error counting users: %w", err)
	}

	// 获取 article count
	articleCount, err := s.articleDao.CountArticles(ctx)
	if err != nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error counting articles: %w", err)
	}

	// 获取 unique views
	uniqueViews, err := s.uniqueViewDao.ListUniqueViews(ctx, startTime, endTime)
	if err != nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error retrieving unique views: %w", err)
	}

	// 获取 article statistics
	articleStatisticsList, err := s.articleDao.ListArticleStatistics(ctx)
	if err != nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error retrieving article statistics: %w", err)
	}

	// 获取 category data
	categoryDTOList, err := s.articleDao.ListCategoryDTO(ctx)
	if err != nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error retrieving category data: %w", err)
	}

	// 获取 tag data
	tagDTOList, err := s.articleDao.ListTagDTO(ctx)
	if err != nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error retrieving tag data: %w", err)
	}

	// 获取 top 5 articles from Redis
	articleMap, err := s.redisService.ZRevRangeWithScore(ctx, "ARTICLE_VIEWS_COUNT", 0, 4)
	if err != nil {
		return dto.BlogBackInfoDTO{}, fmt.Errorf("error retrieving top articles: %w", err)
	}

	// 将 map[interface{}]float64 转换为 map[string]float64
	stringArticleMap := make(map[string]float64)
	for key, score := range articleMap {
		strKey, ok := key.(string)
		if !ok {
			// 处理转换错误
			continue
		}
		stringArticleMap[strKey] = score
	}

	blogBackInfoDTO := dto.BlogBackInfoDTO{
		ViewsCount:            viewsCount,
		MessageCount:          int(messageCount),
		UserCount:             int(userCount),
		ArticleCount:          int(articleCount),
		UniqueViewDTOList:     uniqueViews,
		ArticleStatisticsList: articleStatisticsList,
		CategoryDTOList:       categoryDTOList,
		TagDTOList:            tagDTOList,
	}

	// 获取文章排行
	if len(articleMap) > 0 {
		articleRankDTOList, err := s.listArticleRank(ctx, stringArticleMap)
		if err != nil {
			return dto.BlogBackInfoDTO{}, fmt.Errorf("error retrieving article rank: %w", err)
		}
		blogBackInfoDTO.ArticleRankDTOList = articleRankDTOList
	}

	return blogBackInfoDTO, nil
}

// listArticleRank 根据文章Map获取文章排行
func (s *blogInfoServiceImpl) listArticleRank(ctx context.Context, articleMap map[string]float64) ([]dto.ArticleRankDTO, error) {
	var articleRankList []dto.ArticleRankDTO
	for articleID := range articleMap {
		id, err := strconv.Atoi(articleID)
		if err != nil {
			return nil, fmt.Errorf("error converting article ID: %w", err)
		}
		article, err := s.articleDao.GetArticleById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error retrieving article by ID: %w", err)
		}
		articleRankDTO := dto.ArticleRankDTO{
			ArticleTitle: article.ArticleTitle,
			ViewsCount:   article.ViewsCount,
		}
		articleRankList = append(articleRankList, articleRankDTO)
	}
	return articleRankList, nil
}

func (s *blogInfoServiceImpl) UpdateWebsiteConfig(ctx context.Context, websiteConfigVO vo.WebsiteConfigVO) error {
	log.Println("Updating website configuration")

	// 序列化 websiteConfigVO 为 JSON
	configJSON, err := json.Marshal(websiteConfigVO)
	if err != nil {
		log.Printf("Error marshaling websiteConfigVO to JSON: %v", err)
		return err
	}

	// 创建 WebsiteConfig 指针实例
	websiteConfig := &model.WebsiteConfig{
		ID:     1, // 确保这是正确的 ID
		Config: string(configJSON),
	}

	// 更新数据库中的配置
	err = s.websiteConfigDao.UpdateByID(ctx, s.db, websiteConfig)
	if err != nil {
		log.Printf("Error updating website configuration in database: %v", err)
		return err
	}

	// 删除 Redis 缓存
	_, err = s.redisService.Del(ctx, constants.WebsiteConfig)
	if err != nil {
		log.Printf("Error deleting website configuration from Redis: %v", err)
		return err
	}

	log.Println("Website configuration updated successfully")
	return nil
}

func (s *blogInfoServiceImpl) GetAbout(ctx context.Context) (string, error) {
	value, err := s.redisService.Get(ctx, constants.About)
	if err != nil {
		if err == redis.Nil {
			// Redis 中没有找到，返回空字符串
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (s *blogInfoServiceImpl) UpdateAbout(ctx context.Context, blogInfoVO vo.BlogInfoVO) error {
	return s.redisService.Set(ctx, "ABOUT", blogInfoVO.AboutContent, 0)
}

func (s *blogInfoServiceImpl) Report(ctx context.Context, request *http.Request) {
	// 获取 IP 地址
	ipAddress := utils.GetIPAddress(request)

	// 获取访问设备
	userAgent := utils.GetUserAgent(request)
	browserName, _ := userAgent.Browser() // 获取浏览器名称和版本
	operatingSystem := userAgent.OS()

	// 生成唯一用户标识
	uuid := ipAddress + browserName + operatingSystem
	md5Hash := md5.Sum([]byte(uuid))
	md5 := hex.EncodeToString(md5Hash[:])

	// 判断是否访问
	isMember, err := s.redisService.SIsMember(ctx, "UNIQUE_VISITOR", md5)
	if err != nil {
		log.Printf("Error checking if member exists in Redis: %v", err)
		return
	}

	if !isMember {
		// 统计游客地域分布
		ipSource := utils.GetIPSource(ipAddress)
		if strings.TrimSpace(ipSource) != "" {
			ipSource = strings.Replace(ipSource[:2], "省", "", -1)
			ipSource = strings.Replace(ipSource, "市", "", -1)
			_, err := s.redisService.HIncr(ctx, "VISITOR_AREA", ipSource, 1)
			if err != nil {
				log.Printf("Error incrementing visitor area in Redis: %v", err)
			}
		} else {
			_, err := s.redisService.HIncr(ctx, "VISITOR_AREA", "UNKNOWN", 1)
			if err != nil {
				log.Printf("Error incrementing visitor area in Redis: %v", err)
			}
		}

		// 访问量 +1
		_, err = s.redisService.Incr(ctx, "BLOG_VIEWS_COUNT", 1)
		if err != nil {
			log.Printf("Error incrementing blog views count in Redis: %v", err)
		}

		// 保存唯一标识
		_, err = s.redisService.SAdd(ctx, "UNIQUE_VISITOR", md5)
		if err != nil {
			log.Printf("Error adding unique visitor in Redis: %v", err)
		}
	}
}
