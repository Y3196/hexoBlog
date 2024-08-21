package app

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"goBolg/common"
	"goBolg/config"
	"goBolg/controller"
	"goBolg/dao"
	"goBolg/handler"
	"goBolg/model"
	"goBolg/rabbitmq"
	"goBolg/rabbitmq/rabbitService"
	"goBolg/rabbitmq/rabbitService/rabbitImpl"
	"goBolg/service"
	"goBolg/service/Impl"
	"goBolg/strategy"
	context "goBolg/strategy/contxt"
	"goBolg/strategy/strategyImpl"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// App 包含所有初始化后的实例
type App struct {
	Config      *config.AppConfig
	Database    *gorm.DB
	RedisClient *redis.Client
	RabbitMQ    *amqp.Connection
	Controllers *Controllers
}

// Controllers 包含所有控制器实例
type Controllers struct {
	ArticleController    *controller.ArticleController
	BlogInfoController   *controller.BlogInfoController
	CategoryController   *controller.CategoryController
	CommentController    *controller.CommentController
	FriendLinkController *controller.FriendLinkController
	LogController        *controller.LogController
	MenuController       *controller.MenuController
	MessageController    *controller.MessageController
	PageController       *controller.PageController
	PhotoAlbumController *controller.PhotoAlbumController
	PhotoController      *controller.PhotoController
	ResourceController   *controller.ResourceController
	RoleController       *controller.RoleController
	TagController        *controller.TagController
	TalkController       *controller.TalkController
	UserInfoController   *controller.UserInfoController
	UserAuthController   *controller.UserAuthController
	UserAuthDao          dao.UserAuthDao
	UserInfoDao          dao.UserInfoDao
	RoleDao              dao.RoleDao
	RedisService         service.RedisService
}

// Initialize 初始化应用程序
func Initialize() (*App, error) {
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	// 构建配置文件路径
	configPath := filepath.Join(cwd, "config/application.yaml")

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file '%s' does not exist", configPath)
	}

	// 加载应用程序配置
	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load app config: %w", err)
	}

	// 连接数据库
	database, err := common.ConnectDB(&appConfig.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 连接 Redis
	redisClient, err := common.ConnectRedis(&appConfig.Redis)
	if err != nil {
		common.CloseDB(database)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// 连接 RabbitMQ
	rabbitMQ, err := common.ConnectRabbitMQ(&appConfig.RabbitMQ)
	if err != nil {
		common.CloseDB(database)
		common.CloseRedis(redisClient)
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// 自动迁移数据库结构
	err = database.AutoMigrate(&model.Tag{}, &model.Article{}, &model.Message{}, &model.ArticleTag{})
	if err != nil {
		common.CloseDB(database)
		common.CloseRedis(redisClient)
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// 初始化 RabbitMQ 客户端
	rabbitMQClient, err := rabbitmq.NewRabbitMQClient(fmt.Sprintf("amqp://%s:%s@%s:%d/", appConfig.RabbitMQ.Username, appConfig.RabbitMQ.Password, appConfig.RabbitMQ.Host, appConfig.RabbitMQ.Port))
	if err != nil {
		common.CloseDB(database)
		common.CloseRedis(redisClient)
		common.CloseRabbitMQ(rabbitMQ)
		return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	// 初始化公用的依赖
	articleDao := dao.NewArticleDao(database)
	categoryDao := dao.NewCategoryDao(database)
	tagDao := dao.NewTagDao(database)
	userAuthDao := dao.NewUserAuthDao(database)
	userInfoDao := dao.NewUserInfoDao(database)
	roleDao := dao.NewRoleDao(database)
	redisService := Impl.NewRedisServiceImpl(redisClient)
	websiteConfigDao := dao.NewWebsiteConfigDao(database)
	pageDao := dao.NewPageDao(database)
	//初始化PageService
	pageService := Impl.NewPageService(database, pageDao, redisService)
	messageDao := dao.NewMessageDao(database)
	uniqueViewDao := dao.NewUniqueViewDao(database)

	blogInfoService := Impl.NewBlogInfoService(userInfoDao, messageDao, uniqueViewDao, articleDao, categoryDao, tagDao, redisService, websiteConfigDao, pageService, database)
	categoryService := Impl.NewCategoryService(categoryDao, articleDao)

	// 初始化 CommentService
	commentDao := dao.NewCommentDao(database)
	talkDao := dao.NewTalkDao(database)
	// 配置config的
	commentService := Impl.NewCommentServiceImpl(commentDao, articleDao, talkDao, userInfoDao, redisService, rabbitMQClient, blogInfoService, appConfig.Website.URL)

	//初始化 FriendLinkService
	friendLinkDao := dao.NewFriendLinkDao(database)
	friendLinkService := Impl.NewFriendLinkServiceImpl(friendLinkDao)

	//初始化 OperationLogService
	operationLogDao := dao.NewOperationLogDao(database)
	operationLogService := Impl.NewOperationLogServiceImpl(operationLogDao)

	//初始化  MenuService
	menuDao := dao.NewMenuDao(database)
	roleMenu := dao.NewRoleMenuDao(database)
	menuService := Impl.NewMenuService(menuDao, roleMenu)

	//初始化MessageService
	messageService := Impl.NewMessageService(messageDao, blogInfoService, &http.Request{})

	// 初始化上传策略上下文
	uploadStrategyContext, err := context.NewUploadStrategyContext(appConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize upload strategy contxt: %w", err)
	}

	//初始化 PhotoAlbumService
	photoAlbumDao := dao.NewPhotoAlbumDao(database)
	photoDao := dao.NewPhotoDao(database)
	photoAlbumService := Impl.NewPhotoAlbumService(photoAlbumDao, photoDao)

	// 初始化 PhotoService
	photoService, err := Impl.NewPhotoServiceImpl(photoDao, photoAlbumDao)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize photo rabbitService: %w", err)
	}

	// 初始化 ResourceService
	resourceDao := dao.NewResourceDao(database)
	roleResourceDao := dao.NewRoleResourceDao(database)
	filter := handler.NewFilterInvocationSecurityMetadataSourceImpl(roleDao)
	resourceService := Impl.NewResourceService(resourceDao, roleResourceDao, filter)

	// 初始化 RoleService
	userRoleDao := dao.NewUserRoleDao(database)
	roleMenuDao := dao.NewRoleMenuDao(database)
	roleService := Impl.NewRoleServiceImpl(roleDao, userRoleDao, roleResourceDao, roleMenuDao, database, filter)

	//初始化 TagService
	articleTagDao := dao.NewArticleTagDao(database)
	tagService := Impl.NewTagServiceImpl(tagDao, articleTagDao, database)

	//初始化 TalkService
	talkService := Impl.NewTalkService(talkDao, commentDao, redisService)

	//初始化 UserInfoService
	userInfoService := Impl.NewUserInfoService(userInfoDao, uploadStrategyContext, redisService)

	//初始化 UserAuthService
	rabbitSer := rabbitImpl.NewRabbitService(rabbitMQ)
	emailService := rabbitService.NewEmailService(appConfig.Email)
	log.Printf("Created EmailService with Host: %s, Port: %d", emailService.Host, emailService.Port)
	userAuthService := Impl.NewUserAuthService(*emailService, redisService, rabbitSer, userInfoDao, userRoleDao, userAuthDao, blogInfoService)
	// 初始化控制器
	controllers := &Controllers{
		ArticleController:    NewArticleController(database, redisClient, blogInfoService, appConfig),
		BlogInfoController:   NewBlogInfoController(blogInfoService),
		CategoryController:   NewCategoryController(categoryService),
		CommentController:    NewCommentController(commentService),
		FriendLinkController: NewFriendLinkController(friendLinkService),
		LogController:        NewLogController(operationLogService),
		MenuController:       NewMenuController(menuService),
		MessageController:    NewMessageController(messageService),
		PageController:       NewPagesController(pageService),
		PhotoAlbumController: NewPhotoAlbumController(uploadStrategyContext, photoAlbumService), // 初始化PhotoAlbumController
		PhotoController:      NewPhotoController(photoService),
		ResourceController:   NewResourceController(resourceService),
		RoleController:       NewRoleController(roleService),
		TagController:        NewTagController(tagService),
		TalkController:       NewTalkController(talkService, uploadStrategyContext),
		UserInfoController:   NewUserInfoController(userInfoService),
		UserAuthController:   NewUserAuthController(userAuthService),
		UserAuthDao:          userAuthDao,
		UserInfoDao:          userInfoDao,
		RoleDao:              roleDao,
		RedisService:         redisService,
	}

	// 创建 App 实例
	app := &App{
		Config:      appConfig,
		Database:    database,
		RedisClient: redisClient,
		RabbitMQ:    rabbitMQ,
		Controllers: controllers,
	}

	return app, nil
}

// NewArticleController 初始化文章控制器
func NewArticleController(database *gorm.DB, redisClient *redis.Client, blogInfoService service.BlogInfoService, appConfig *config.AppConfig) *controller.ArticleController {
	articleDao := dao.NewArticleDao(database)
	articleTagDao := dao.NewArticleTagDao(database)
	categoryDao := dao.NewCategoryDao(database)
	tagDao := dao.NewTagDao(database)
	tagService := Impl.NewTagServiceImpl(tagDao, articleTagDao, database)
	redisService := Impl.NewRedisServiceImpl(redisClient)

	// 初始化搜索策略
	searchStrategyMap := map[string]strategy.SearchStrategy{
		"mysql": &strategyImpl.MySqlSearchStrategyImpl{ArticleDao: articleDao},
	}
	searchStrategyContext := context.NewSearchStrategyContext(appConfig.Search.Mode, searchStrategyMap)

	articleService := Impl.NewArticleServiceImpl(articleDao, articleTagDao, categoryDao, tagDao, tagService, redisService, blogInfoService, searchStrategyContext, database)
	return &controller.ArticleController{

		Service: articleService,
	}
}

// NewBlogInfoController 初始化博客信息控制器
func NewBlogInfoController(blogInfoService service.BlogInfoService) *controller.BlogInfoController {
	return &controller.BlogInfoController{
		BlogInfoService: blogInfoService,
	}
}

// NewCategoryController 初始化分类控制器
func NewCategoryController(categoryService service.CategoryService) *controller.CategoryController {
	return &controller.CategoryController{
		CategoryService: categoryService,
	}
}

// NewCommentController 初始化评论控制器
func NewCommentController(commentService service.CommentService) *controller.CommentController {
	return &controller.CommentController{
		CommentService: commentService,
	}
}

// NewFriendLinkController 初始化友链控制器
func NewFriendLinkController(friendLinkService service.FriendLinkService) *controller.FriendLinkController {
	return &controller.FriendLinkController{
		FriendLinkService: friendLinkService,
	}
}

// NewFriendLinkController 初始化友链控制器
func NewLogController(operationService service.OperationLogService) *controller.LogController {
	return &controller.LogController{
		OperationLogService: operationService,
	}
}

// NewFriendLinkController 初始化友链控制器
func NewMenuController(menuService service.MenuService) *controller.MenuController {
	return &controller.MenuController{
		MenuService: menuService,
	}
}

// NewMessageController 初始化留言板控制器
func NewMessageController(messageService service.MessageService) *controller.MessageController {
	return &controller.MessageController{
		MessageService: messageService,
	}
}

// 初始化页面控制器
func NewPagesController(pageService service.PageService) *controller.PageController {
	return &controller.PageController{
		PageService: pageService,
	}
}

// NewPhotoAlbumController 初始化相册控制器
func NewPhotoAlbumController(uploadStrategyContext *context.UploadStrategyContext, albumService service.PhotoAlbumService) *controller.PhotoAlbumController {
	return &controller.PhotoAlbumController{
		UploadStrategyContext: uploadStrategyContext,
		PhotoAlbumService:     albumService,
	}
}

// NewPhotoController 初始化照片控制器
func NewPhotoController(photoService service.PhotoService) *controller.PhotoController {
	return &controller.PhotoController{
		PhotoService: photoService,
	}
}

// NewResourceController 初始化资源控制器
func NewResourceController(resourceService service.ResourceService) *controller.ResourceController {
	return &controller.ResourceController{
		ResourceService: resourceService,
	}
}

// NewRoleController 初始化角色控制器
func NewRoleController(roleService service.RoleService) *controller.RoleController {
	return &controller.RoleController{
		RoleService: roleService,
	}
}

// NewTagController 初始化标签控制器
func NewTagController(tagService service.TagService) *controller.TagController {
	return &controller.TagController{
		Service: tagService,
	}
}

// NewTalkController 初始化说说控制器
func NewTalkController(talkService service.TalkService, uploadStrategyContext *context.UploadStrategyContext) *controller.TalkController {
	return &controller.TalkController{
		UploadStrategyContext: uploadStrategyContext,
		TalkService:           talkService,
	}
}

// NewUserInfoController 初始化用户信息控制器
func NewUserInfoController(userInfoService service.UserInfoService) *controller.UserInfoController {
	return &controller.UserInfoController{
		userInfoService,
	}
}

// NewUserAuthController 初始化用户认证控制器
func NewUserAuthController(userAuthService service.UserAuthService) *controller.UserAuthController {
	return &controller.UserAuthController{
		userAuthService,
	}
}

// Close 关闭应用程序
func (app *App) Close() {
	common.CloseDB(app.Database)
	common.CloseRedis(app.RedisClient)
	common.CloseRabbitMQ(app.RabbitMQ)
}
