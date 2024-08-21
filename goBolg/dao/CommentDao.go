package dao

import (
	"context"
	"fmt"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
	"log"
)

// CommentDao 评论 DAO 接口
type CommentDao interface {
	ListComments(ctx context.Context, current int, size int, commentVO vo.CommentVO) ([]dto.CommentDTO, error)
	ListReplies(ctx context.Context, commentIdList []int) ([]dto.ReplyDTO, error)
	ListRepliesByCommentId(ctx context.Context, current int, size int, commentId int) ([]dto.ReplyDTO, error)
	ListReplyCountByCommentId(ctx context.Context, commentIdList []int) ([]dto.ReplyCountDTO, error)
	ListCommentCountByTopicIds(ctx context.Context, topicIdList []int) ([]dto.CommentCountDTO, error)
	ListCommentBackDTO(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]dto.CommentBackDTO, error)
	CountCommentDTO(ctx context.Context, condition vo.ConditionVO) (int, error)
	Insert(ctx context.Context, comment model.Comment) (int, error)
	UpdateBatchByID(ctx context.Context, comments []model.Comment) error
	RemoveByIds(ctx context.Context, ids []int) error

	GetCommentByID(ctx context.Context, commentID int) (dto.CommentDTO, error)
}

type commentDao struct {
	db *gorm.DB
}

// NewCommentDao 创建新的 CommentDao 实例
func NewCommentDao(db *gorm.DB) CommentDao {
	return &commentDao{db: db}
}

// ListComments 查询评论
func (dao *commentDao) ListComments(ctx context.Context, current int, size int, commentVO vo.CommentVO) ([]dto.CommentDTO, error) {
	var comments []dto.CommentDTO

	offset := (current - 1) * size
	query := `
        SELECT
            u.nickname,
            u.avatar,
            u.web_site,
            c.user_id,
            c.id,
            c.comment_content,
            c.create_time
        FROM
            tb_comment c
            JOIN tb_user_info u ON c.user_id = u.id
        WHERE
            c.is_review = 1
            AND c.parent_id IS NULL
            AND (c.topic_id = ? OR ? IS NULL)
            AND c.type = ?
        ORDER BY
            c.create_time DESC -- 这里改为按创建时间降序排列
        LIMIT ?, ?`

	// 强制查询主库，如果使用了读写分离
	err := dao.db.WithContext(ctx).Raw(query, commentVO.TopicID, commentVO.TopicID, commentVO.Type, offset, size).Scan(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// 返回空切片而不是 nil，避免空指针异常
	if len(comments) == 0 {
		return []dto.CommentDTO{}, nil
	}

	return comments, nil
}

// ListReplies 查询评论id集合下的回复
func (dao *commentDao) ListReplies(ctx context.Context, commentIdList []int) ([]dto.ReplyDTO, error) {
	var replies []dto.ReplyDTO

	query := `
		SELECT
			c.user_id,
			u.nickname,
			u.avatar,
			u.web_site,
			c.reply_user_id,
			r.nickname AS reply_nickname,
			r.web_site AS reply_web_site,
			c.id,
			c.parent_id,
			c.comment_content,
			c.create_time
		FROM
			tb_comment c
			JOIN tb_user_info u ON c.user_id = u.id
			JOIN tb_user_info r ON c.reply_user_id = r.id
		WHERE
			c.is_review = 1
			AND c.parent_id IN ?`

	err := dao.db.WithContext(ctx).Raw(query, commentIdList).Scan(&replies).Error
	if err != nil {
		return nil, err
	}

	return replies, nil
}

// ListRepliesByCommentId 查询某条评论下的回复
func (dao *commentDao) ListRepliesByCommentId(ctx context.Context, current int, size int, commentId int) ([]dto.ReplyDTO, error) {
	var replies []dto.ReplyDTO

	offset := (current - 1) * size
	query := `
		SELECT
			c.user_id,
			u.nickname,
			u.avatar,
			u.web_site,
			c.reply_user_id,
			r.nickname AS reply_nickname,
			r.web_site AS reply_web_site,
			c.id,
			c.parent_id,
			c.comment_content,
			c.create_time
		FROM
			tb_comment c
			JOIN tb_user_info u ON c.user_id = u.id
			JOIN tb_user_info r ON c.reply_user_id = r.id
		WHERE
			c.is_review = 1
			AND c.parent_id = ?
		ORDER BY
			c.create_time ASC
		LIMIT ?, ?`

	err := dao.db.WithContext(ctx).Raw(query, commentId, offset, size).Scan(&replies).Error
	if err != nil {
		return nil, err
	}

	return replies, nil
}

// ListReplyCountByCommentId 查询评论id集合下的回复总量
func (dao *commentDao) ListReplyCountByCommentId(ctx context.Context, commentIdList []int) ([]dto.ReplyCountDTO, error) {
	var replyCounts []dto.ReplyCountDTO

	query := `
		SELECT
			parent_id AS comment_id,
			COUNT(1) AS reply_count
		FROM
			tb_comment
		WHERE
			is_review = 1
			AND parent_id IN ?
		GROUP BY
			parent_id`

	err := dao.db.WithContext(ctx).Raw(query, commentIdList).Scan(&replyCounts).Error
	if err != nil {
		return nil, err
	}

	return replyCounts, nil
}

// ListCommentCountByTopicIds 查询评论量
func (dao *commentDao) ListCommentCountByTopicIds(ctx context.Context, topicIdList []int) ([]dto.CommentCountDTO, error) {
	var commentCounts []dto.CommentCountDTO

	query := `
		SELECT
			topic_id AS id,
			COUNT(1) AS comment_count
		FROM
			tb_comment
		WHERE
			topic_id IN ?
			AND parent_id IS NULL
		GROUP BY
			topic_id`

	err := dao.db.WithContext(ctx).Raw(query, topicIdList).Scan(&commentCounts).Error
	if err != nil {
		return nil, err
	}

	return commentCounts, nil
}

// ListCommentBackDTO 查询后台评论
func (dao *commentDao) ListCommentBackDTO(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]dto.CommentBackDTO, error) {
	var comments []dto.CommentBackDTO

	offset := (current - 1) * size
	query := dao.db.Table("tb_comment").
		Select(`
			tb_comment.id,
			u.avatar,
			u.nickname,
			r.nickname AS reply_nickname,
			a.article_title,
			tb_comment.comment_content,
			tb_comment.type,
			tb_comment.create_time,
			tb_comment.is_review
		`).
		Joins("LEFT JOIN tb_article a ON tb_comment.topic_id = a.id").
		Joins("LEFT JOIN tb_user_info u ON tb_comment.user_id = u.id").
		Joins("LEFT JOIN tb_user_info r ON tb_comment.reply_user_id = r.id").
		Where("1 = 1")

	if condition.Type != nil {
		query = query.Where("tb_comment.type = ?", *condition.Type)
	}
	if condition.IsReview != nil {
		query = query.Where("tb_comment.is_review = ?", *condition.IsReview)
	}
	if condition.Keywords != nil {
		query = query.Where("u.nickname LIKE ?", "%"+*condition.Keywords+"%")
	}
	log.Printf("Executing list query: %v", query)
	err := query.Order("tb_comment.id DESC").
		Offset(offset).
		Limit(size).
		Scan(&comments).Error
	if err != nil {
		log.Printf("Error executing list query: %v", err)
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	return comments, nil
}

// CountCommentDTO 统计后台评论数量
func (dao *commentDao) CountCommentDTO(ctx context.Context, condition vo.ConditionVO) (int, error) {
	var count int64
	query := dao.db.Table("tb_comment").
		Joins("LEFT JOIN tb_user_info u ON tb_comment.user_id = u.id").
		Where("1 = 1")

	if condition.Type != nil {
		query = query.Where("tb_comment.type = ?", *condition.Type)
	}
	if condition.IsReview != nil {
		query = query.Where("tb_comment.is_review = ?", *condition.IsReview)
	}
	if condition.Keywords != nil {
		query = query.Where("u.nickname LIKE ?", "%"+*condition.Keywords+"%")
	}
	log.Printf("Executing count query: %v", query)
	err := query.Count(&count).Error
	if err != nil {
		log.Printf("Error executing count query: %v", err)
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return int(count), nil
}

func (dao *commentDao) Insert(ctx context.Context, comment model.Comment) (int, error) {
	result := dao.db.WithContext(ctx).Create(&comment)
	if result.Error != nil {
		return 0, result.Error
	}
	// Return the generated ID
	return comment.ID, nil
}

func (dao *commentDao) UpdateBatchByID(ctx context.Context, comments []model.Comment) error {
	return dao.db.WithContext(ctx).Model(&model.Comment{}).Updates(comments).Error
}

// RemoveByIds 批量删除评论
func (dao *commentDao) RemoveByIds(ctx context.Context, ids []int) error {
	if err := dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Comment{}).Error; err != nil {
		return err
	}
	return nil
}

func (dao *commentDao) GetCommentByID(ctx context.Context, commentID int) (dto.CommentDTO, error) {
	var comment dto.CommentDTO
	query := `
        SELECT
            u.nickname,
            u.avatar,
            u.web_site,
            c.user_id,
            c.id,
            c.comment_content,
            c.create_time
        FROM
            tb_comment c
            JOIN tb_user_info u ON c.user_id = u.id
        WHERE
            c.id = ?`

	err := dao.db.WithContext(ctx).Debug().Raw(query, commentID).Scan(&comment).Error
	if err != nil {
		return dto.CommentDTO{}, fmt.Errorf("failed to retrieve comment: %w", err)
	}

	return comment, nil
}
