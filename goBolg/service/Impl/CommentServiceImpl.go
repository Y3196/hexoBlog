package Impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	constants "goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/enums"
	"goBolg/model"
	"goBolg/rabbitmq"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"golang.org/x/sync/errgroup"
	"log"
	"strconv"
	"time"
)

// commentServiceImpl 实现了 CommentService 接口
type commentServiceImpl struct {
	commentDao      dao.CommentDao
	articleDao      dao.ArticleDao
	talkDao         dao.TalkDao
	userInfoDao     dao.UserInfoDao
	rabbitMQClient  *rabbitmq.RabbitMQClient
	redisService    service.RedisService
	blogInfoService service.BlogInfoService
	websiteUrl      string
}

// NewCommentServiceImpl 创建一个新的 commentServiceImpl 实例
func NewCommentServiceImpl(commentDao dao.CommentDao, articleDao dao.ArticleDao, talkDao dao.TalkDao, userInfoDao dao.UserInfoDao, redisService service.RedisService, rabbitMQClient *rabbitmq.RabbitMQClient, blogInfoService service.BlogInfoService, websiteUrl string) service.CommentService {
	return &commentServiceImpl{
		commentDao:      commentDao,
		articleDao:      articleDao,
		talkDao:         talkDao,
		userInfoDao:     userInfoDao,
		redisService:    redisService,
		rabbitMQClient:  rabbitMQClient,
		blogInfoService: blogInfoService,
		websiteUrl:      websiteUrl,
	}
}

// ListComments 实现了 CommentService 接口中的 ListComments 方法
func (s *commentServiceImpl) ListComments(ctx context.Context, commentVO vo.CommentVO, current int, size int) ([]dto.CommentDTO, int, error) {
	commentDTOList, err := s.commentDao.ListComments(ctx, current, size, commentVO)
	if err != nil {
		return nil, 0, err
	}

	if len(commentDTOList) == 0 {
		return []dto.CommentDTO{}, 0, nil
	}

	// 查询 Redis 的评论点赞数据
	likeCountMap, err := s.redisService.HGetAll(ctx, constants.CommentLikeCount)
	if err != nil {
		log.Printf("Error getting like count from Redis: %v", err)
		return nil, 0, err
	}

	// 提取评论 ID 集合
	var commentIdList []int
	for _, comment := range commentDTOList {
		commentIdList = append(commentIdList, comment.ID)
	}

	// 根据评论 ID 集合查询回复数据
	replyDTOList, err := s.commentDao.ListReplies(ctx, commentIdList)
	if err != nil {
		return nil, 0, err
	}

	// 封装回复点赞量
	for i := range replyDTOList {
		replyIDStr := strconv.Itoa(replyDTOList[i].ID)
		if likeCount, exists := likeCountMap[replyIDStr]; exists {
			replyDTOList[i].LikeCount, _ = strconv.Atoi(likeCount)
		}
	}

	// 根据评论 ID 分组回复数据
	replyMap := make(map[int][]dto.ReplyDTO)
	for _, reply := range replyDTOList {
		replyMap[reply.ParentID] = append(replyMap[reply.ParentID], reply)
	}

	// 根据评论 ID 查询回复量
	replyCountMap, err := s.commentDao.ListReplyCountByCommentId(ctx, commentIdList)
	if err != nil {
		return nil, 0, err
	}

	// 封装评论数据
	for i := range commentDTOList {
		commentIDStr := strconv.Itoa(commentDTOList[i].ID)
		if likeCount, exists := likeCountMap[commentIDStr]; exists {
			commentDTOList[i].LikeCount, _ = strconv.Atoi(likeCount)
		}
		commentDTOList[i].ReplyDTOList = replyMap[commentDTOList[i].ID]

		for _, count := range replyCountMap {
			if count.CommentID == commentDTOList[i].ID {
				commentDTOList[i].ReplyCount = count.ReplyCount
			}
		}
	}

	return commentDTOList, len(commentDTOList), nil
}

// ListRepliesByCommentId 实现了 CommentService 接口中的 ListRepliesByCommentId 方法
func (s *commentServiceImpl) ListRepliesByCommentId(ctx context.Context, commentId int) ([]dto.ReplyDTO, error) {
	current := utils.GetCurrent(ctx)
	size := utils.GetSize(ctx)

	// Log the pagination parameters and comment ID
	log.Printf("Fetching replies for commentId: %d, current page: %d, size: %d", commentId, current, size)

	replyDTOList, err := s.commentDao.ListRepliesByCommentId(ctx, current, size, commentId)
	if err != nil {
		log.Printf("Database query failed for comment ID: %d, error: %v", commentId, err)
		return nil, err
	}

	likeCountMap, err := s.redisService.HGetAll(ctx, constants.CommentLikeCount)
	if err != nil {
		log.Printf("Error getting like count from Redis for comment ID: %d, error: %v", commentId, err)
		return nil, err
	}

	for i := range replyDTOList {
		replyIDStr := strconv.Itoa(replyDTOList[i].ID)
		if likeCount, ok := likeCountMap[replyIDStr]; ok {
			replyDTOList[i].LikeCount, _ = strconv.Atoi(likeCount)
		} else {
			replyDTOList[i].LikeCount = 0
			log.Printf("No like count found for reply ID: %d", replyDTOList[i].ID)
		}
	}

	log.Printf("Successfully fetched replies for comment ID: %d", commentId)
	return replyDTOList, nil
}

// SaveComment saves a new comment.
func (s *commentServiceImpl) SaveComment(ctx context.Context, commentVO vo.CommentVO) (dto.CommentDTO, error) {
	// 获取网站配置
	log.Println("获取网站配置")
	websiteConfig, err := s.blogInfoService.GetWebsiteConfig(ctx)
	if err != nil {
		log.Printf("获取网站配置失败: %v", err)
		return dto.CommentDTO{}, fmt.Errorf("failed to get website config: %w", err)
	}

	// 判断是否需要审核
	isReview := websiteConfig.IsCommentReview
	// 过滤标签
	htmlUtils := utils.HTMLUtils{}
	commentContent := htmlUtils.Filter(commentVO.CommentContent)
	log.Printf("过滤后的评论内容: %s", commentContent)

	user, ok := utils.GetLoginUser(ctx)
	if !ok {
		log.Println("用户未登录")
		return dto.CommentDTO{}, fmt.Errorf("user not logged in")
	}
	log.Printf("登录用户: %+v\n", user)

	var replyUserID, parentID *int
	if commentVO.ReplyUserID != 0 {
		replyUserID = &commentVO.ReplyUserID
	}
	if commentVO.ParentID != 0 {
		parentID = &commentVO.ParentID
	}

	comment := model.Comment{
		UserID:         user.UserInfoID,
		ReplyUserID:    replyUserID,
		TopicID:        commentVO.TopicID,
		CommentContent: commentContent,
		ParentID:       parentID,
		Type:           commentVO.Type,
		IsReview:       constants.False,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}

	if isReview == constants.False {
		comment.IsReview = constants.True
	}

	// 在插入前打印
	log.Printf("Before insert - Comment: %+v", comment)
	insertedID, err := s.commentDao.Insert(ctx, comment)
	if err != nil {
		log.Printf("评论插入失败: %v", err)
		return dto.CommentDTO{}, fmt.Errorf("failed to insert comment: %w", err)
	}

	// 在插入后打印
	log.Printf("After insert - Comment: %+v", comment)

	// 获取插入后的评论数据
	savedComment, err := s.commentDao.GetCommentByID(ctx, insertedID)
	if err != nil {
		log.Printf("获取插入后的评论失败: %v", err)
		return dto.CommentDTO{}, fmt.Errorf("failed to retrieve saved comment: %w", err)
	}
	log.Printf("获取插入后的评论数据是: %+v\n", savedComment)

	// 判断是否开启邮箱通知, 通知用户
	if websiteConfig.IsEmailNotice == constants.True {
		log.Println("发送邮件通知")
		err = s.sendNotification(ctx, comment)
		if err != nil {
			log.Printf("发送通知失败，但评论已保存: %v", err)
			// 此处记录通知发送失败的日志，但不返回错误，确保评论保存的事务独立完成。
		}
	}

	return savedComment, nil
}

func (s *commentServiceImpl) sendNotification(ctx context.Context, comment model.Comment) error {
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.Notice(ctx, comment)
	})
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	return nil
}

// Notice 通知评论用户
func (s *commentServiceImpl) Notice(ctx context.Context, comment model.Comment) error {
	// 初始化 userID 为博主ID作为默认值
	userID := constants.BloggerID
	log.Println("Initial userID set to BloggerID:", userID)

	// 如果有回复用户ID，优先使用
	if comment.ReplyUserID != nil && *comment.ReplyUserID != 0 {
		userID = *comment.ReplyUserID
		log.Println("ReplyUserID found, userID set to:", userID)
	} else {
		// 根据评论类型获取对应的用户ID
		commentEnum := enums.GetCommentEnum(comment.Type)
		if commentEnum == nil {
			return fmt.Errorf("unknown comment type: %d", comment.Type)
		}

		switch commentEnum.Type {
		case enums.ARTICLE.Type:
			article, err := s.articleDao.GetArticleById(ctx, comment.TopicID)
			if err != nil {
				return fmt.Errorf("failed to get article by ID: %w", err)
			}
			userID = article.UserID
			log.Println("Article found, userID set to:", userID)

		case enums.TALK.Type:
			talk, err := s.talkDao.GetTalkById(ctx, comment.TopicID)
			if err != nil {
				return fmt.Errorf("failed to get talk by ID: %w", err)
			}
			userID = talk.UserID
			log.Println("Talk found, userID set to:", userID)

		default:
			return fmt.Errorf("unsupported comment type: %d", comment.Type)
		}
	}

	// 检查 userID 是否有效
	if userID == 0 {
		log.Println("Error: userID is 0, aborting Notice method")
		return fmt.Errorf("invalid userID: %d", userID)
	}

	// 获取用户信息
	user, err := s.userInfoDao.GetUserInfoById(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user info by ID: %w", err)
	}

	email := user.Email
	log.Println("User email found:", email)

	if email != "" {
		// 发送消息
		emailDTO := dto.EmailDTO{}
		if comment.IsReview == constants.True {
			// 评论提醒
			emailDTO.Email = email
			emailDTO.Subject = "评论提醒"
			url := s.websiteUrl + enums.GetCommentPath(comment.Type) + strconv.Itoa(comment.TopicID)
			emailDTO.Content = fmt.Sprintf("您收到了一条新的回复，请前往 %s 页面查看", url)
			log.Println("Preparing to send comment notification email to:", email)
		} else {
			// 管理员审核提醒
			adminUser, err := s.userInfoDao.GetUserInfoById(ctx, constants.BloggerID)
			if err != nil {
				return fmt.Errorf("failed to get admin user info: %w", err)
			}
			emailDTO.Email = adminUser.Email
			emailDTO.Subject = "审核提醒"
			emailDTO.Content = "您收到了一条新的回复，请前往后台管理页面审核"
			log.Println("Preparing to send review notification email to:", adminUser.Email)
		}

		messageBytes, err := json.Marshal(emailDTO)
		if err != nil {
			return fmt.Errorf("failed to marshal emailDTO: %w", err)
		}

		err = s.rabbitMQClient.Publish("EMAIL_EXCHANGE", "*", messageBytes)
		if err != nil {
			return fmt.Errorf("failed to send message to RabbitMQ: %w", err)
		}
		log.Println("Email notification sent successfully")
	} else {
		// 如果邮箱为空，记录日志或返回错误
		log.Println("No email found for userID:", userID)
		return fmt.Errorf("email is empty for userID: %d", userID)
	}

	return nil
}

func (s *commentServiceImpl) SaveCommentLike(ctx context.Context, commentId int) (int, error) {
	// 获取当前用户 ID
	user, ok := utils.GetLoginUser(ctx)
	if !ok {
		return 0, fmt.Errorf("user not logged in")
	}
	userId := user.UserInfoID

	commentLikeKey := constants.CommentUserLike + strconv.Itoa(userId)
	commentIdStr := strconv.Itoa(commentId)

	log.Printf("User %d is attempting to like/unlike comment %d", userId, commentId)

	// 判断是否点赞
	isMember, err := s.redisService.SIsMember(ctx, commentLikeKey, commentIdStr)
	if err != nil {
		log.Printf("Error checking membership in set: %v", err)
		return 0, fmt.Errorf("failed to check if member exists: %w", err)
	}

	log.Printf("Is comment %d liked by user %d? %v", commentId, userId, isMember)

	// 获取当前点赞数量
	likeCountStr, err := s.redisService.HGet(ctx, constants.CommentLikeCount, commentIdStr)
	if err != nil && err != redis.Nil {
		log.Printf("Error getting like count for comment %d: %v", commentId, err)
		return 0, fmt.Errorf("failed to get like count: %w", err)
	}

	if likeCountStr == "" || likeCountStr == "nil" {
		log.Printf("Like count for comment %d is not set, initializing to 0", commentId)
		likeCountStr = "0"
		_, err := s.redisService.HSet(ctx, constants.CommentLikeCount, commentIdStr, "0")
		if err != nil {
			log.Printf("Error initializing like count for comment %d: %v", commentId, err)
			return 0, fmt.Errorf("failed to initialize like count: %w", err)
		}
	}

	likeCount, err := strconv.Atoi(likeCountStr)
	if err != nil {
		log.Printf("Error converting like count to int for comment %d: %v", commentId, err)
		return 0, fmt.Errorf("failed to convert like count to int: %w", err)
	}

	// 防止负值问题
	if likeCount < 0 {
		log.Printf("Like count for comment %d is negative (%d), resetting to 0", commentId, likeCount)
		likeCount = 0
		_, err := s.redisService.HSet(ctx, constants.CommentLikeCount, commentIdStr, "0")
		if err != nil {
			log.Printf("Error resetting like count for comment %d: %v", commentId, err)
			return 0, fmt.Errorf("failed to reset like count: %w", err)
		}
	}

	log.Printf("Current like count for comment %d is %d", commentId, likeCount)

	if isMember {
		// 用户已经点赞，现在要取消点赞
		removed, err := s.redisService.SRemove(ctx, commentLikeKey, commentIdStr)
		if err != nil {
			log.Printf("Error removing comment %d from user %d's liked set: %v", commentId, userId, err)
			return 0, fmt.Errorf("failed to remove comment id from set: %w", err)
		}
		log.Printf("Removed from set: %d", removed)

		if likeCount > 0 {
			_, err = s.redisService.HDecr(ctx, constants.CommentLikeCount, commentIdStr, 1)
			if err != nil {
				log.Printf("Error decrementing like count for comment %d: %v", commentId, err)
				return 0, fmt.Errorf("failed to decrement comment like count: %w", err)
			}
			log.Printf("Decremented like count for comment %d by user %d", commentId, userId)
		} else {
			log.Printf("Like count for comment %d is already 0 or negative, not decrementing", commentId)
		}
	} else {
		// 用户没有点赞，现在要点赞
		added, err := s.redisService.SAdd(ctx, commentLikeKey, commentIdStr)
		if err != nil {
			log.Printf("Error adding comment %d to user %d's liked set: %v", commentId, userId, err)
			return 0, fmt.Errorf("failed to add comment id to set: %w", err)
		}
		log.Printf("Added to set: %d", added)

		// 增加点赞数量
		incr, err := s.redisService.HIncr(ctx, constants.CommentLikeCount, commentIdStr, 1)
		if err != nil {
			log.Printf("Error incrementing like count for comment %d: %v", commentId, err)
			return 0, fmt.Errorf("failed to increment comment like count: %w", err)
		}
		log.Printf("Incremented like count for comment %d by user %d, new count: %d", commentId, userId, incr)
	}

	// 获取更新后的点赞数量
	likeCountStr, err = s.redisService.HGet(ctx, constants.CommentLikeCount, commentIdStr)
	if err != nil {
		log.Printf("Error getting updated like count for comment %d: %v", commentId, err)
		return 0, fmt.Errorf("failed to get updated like count: %w", err)
	}

	likeCount, err = strconv.Atoi(likeCountStr)
	if err != nil {
		log.Printf("Error converting updated like count to int for comment %d: %v", commentId, err)
		return 0, fmt.Errorf("failed to convert updated like count to int: %w", err)
	}

	log.Printf("Final like count for comment %d is %d", commentId, likeCount)

	return likeCount, nil
}

func (s *commentServiceImpl) UpdateCommentsReview(ctx context.Context, reviewVO vo.ReviewVO) error {
	comments := make([]model.Comment, len(reviewVO.IDList))
	for i, id := range reviewVO.IDList {
		comments[i] = model.Comment{
			ID:       id,
			IsReview: reviewVO.IsReview,
		}
	}
	return s.commentDao.UpdateBatchByID(ctx, comments)
}

func (s *commentServiceImpl) ListCommentBackDTO(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.CommentBackDTO], error) {
	// 获取当前分页参数
	current := utils.GetCurrent(ctx)
	size := utils.GetSize(ctx)

	// 统计后台评论量
	count, err := s.commentDao.CountCommentDTO(ctx, condition)
	if err != nil {
		log.Printf("Error counting comments: %v", err)
		return vo.PageResult[dto.CommentBackDTO]{}, fmt.Errorf("failed to count comments: %w", err)
	}
	if count == 0 {
		log.Printf("No comments found")
		return vo.PageResult[dto.CommentBackDTO]{}, nil
	}

	// 查询后台评论集合
	comments, err := s.commentDao.ListCommentBackDTO(ctx, current, size, condition)
	if err != nil {
		log.Printf("Error listing comments: %v", err)
		return vo.PageResult[dto.CommentBackDTO]{}, fmt.Errorf("failed to list comments: %w", err)
	}

	return vo.NewPageResult(comments, count), nil
}

func (s *commentServiceImpl) RemoveComments(ctx context.Context, commentIdList []int) error {
	if err := s.commentDao.RemoveByIds(ctx, commentIdList); err != nil {
		log.Printf("Error removing comments: %v", err)
		return err
	}
	return nil
}
