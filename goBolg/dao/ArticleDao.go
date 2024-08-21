package dao

import (
	"context"
	"fmt"
	"goBolg/dto"
	"goBolg/enums"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

// ArticleDao 文章DAO接口
type ArticleDao interface {
	ListArchives(offset, limit int) ([]dto.ArchiveDTO, error)

	CountArchives() (int, error)

	ListArticles(current int, size int) ([]dto.ArticleHomeDTO, error)

	CountArticles(ctx context.Context) (int64, error)

	ListArticleStatistics(ctx context.Context) ([]dto.ArticleStatisticsDTO, error)

	GetArticleById(ctx context.Context, articleId int) (dto.ArticleDTO, error)

	ListArticleBacks(ctx context.Context, current, size int, condition vo.ConditionVO) ([]dto.ArticleBackDTO, error)

	CountArticleBacks(ctx context.Context, condition vo.ConditionVO) (int, error)

	ListRecommendArticles(ctx context.Context, articleId int) ([]dto.ArticleRecommendDTO, error)

	ListCategoryDTO(ctx context.Context) ([]dto.CategoryDTO, error)

	ListTagDTO(ctx context.Context) ([]dto.TagDTO, error)

	// 查询上一篇文章
	GetPreviousArticle(ctx context.Context, articleId int) (model.Article, error)

	// 查询下一篇文章
	GetNextArticle(ctx context.Context, articleId int) (model.Article, error)

	// 根据条件查询文章
	ListArticlesByCondition(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]dto.ArticlePreviewDTO, error)

	SaveOrUpdate(ctx context.Context, article *model.Article) error

	UpdateById(ctx context.Context, article *model.Article) error

	UpdateBatchById(ctx context.Context, articles []model.Article) error

	DeleteArticlesByIds(ctx context.Context, articleIdList []int) error

	GetArticleCountByCategoryIDs(ctx context.Context, categoryIDList []int) (int64, error)

	GetDb() *gorm.DB
}

type articleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao { // 注意这里返回的是接口类型 ArticleDao
	return &articleDao{db: db}
}

// ListArchives 查询文章归档列表并支持分页
func (r *articleDao) ListArchives(offset, limit int) ([]dto.ArchiveDTO, error) {
	var archives []dto.ArchiveDTO

	result := r.db.Table("tb_article").
		Select("id, article_title, create_time").
		Offset(offset).
		Limit(limit).
		Order("create_time DESC").
		Find(&archives)

	if result.Error != nil {
		log.Printf("Error executing query: %v", result.Error)
		return nil, result.Error
	}

	log.Printf("Query result: %+v", archives)
	return archives, nil
}

func (r *articleDao) CountArchives() (int, error) {
	var count int64

	result := r.db.Table("tb_article").Count(&count)
	if result.Error != nil {
		log.Printf("Error executing count query: %v", result.Error)
		return 0, result.Error
	}

	return int(count), nil
}

func (dao *articleDao) ListArticles(offset int, size int) ([]dto.ArticleHomeDTO, error) {
	var articles []dto.ArticleHomeDTO
	query := `
    SELECT
        a.id,
        a.article_cover,
        a.article_title,
        SUBSTR(a.article_content, 1, 500) AS article_content,
        a.create_time,
        a.type,
        a.is_top,
        a.category_id,
        c.category_name,
        GROUP_CONCAT(DISTINCT t.id) AS tag_ids,
        GROUP_CONCAT(DISTINCT t.tag_name) AS tag_names
    FROM
        tb_article a
    JOIN tb_category c ON a.category_id = c.id
    LEFT JOIN tb_article_tag atg ON a.id = atg.article_id
    LEFT JOIN tb_tag t ON t.id = atg.tag_id
    WHERE
        a.is_delete = 0
        AND a.status = 1
    GROUP BY
        a.id, a.article_cover, a.article_title, a.article_content, a.create_time, a.type, a.is_top, a.category_id, c.category_name
    ORDER BY
        a.is_top DESC,
        a.id DESC
    LIMIT ?, ?`

	rows, err := dao.db.Raw(query, offset, size).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articleMap := make(map[int]*dto.ArticleHomeDTO)
	for rows.Next() {
		var article dto.ArticleHomeDTO
		var tagIDs string
		var tagNames string
		if err := rows.Scan(&article.ID, &article.ArticleCover, &article.ArticleTitle, &article.ArticleContent, &article.CreateTime, &article.Type,
			&article.IsTop, &article.CategoryID, &article.CategoryName, &tagIDs, &tagNames); err != nil {
			return nil, err
		}

		// Parse tags
		tagIDList := strings.Split(tagIDs, ",")
		tagNameList := strings.Split(tagNames, ",")

		tags := make([]dto.TagDTO, len(tagIDList))
		for i := range tagIDList {
			tagID, err := strconv.Atoi(tagIDList[i])
			if err != nil {
				return nil, err // 错误处理：无法将标签 ID 转换为整数
			}
			tags[i] = dto.TagDTO{
				ID:      tagID,
				TagName: tagNameList[i],
			}
		}
		article.Tags = tags

		// Ensure unique articles
		if existingArticle, found := articleMap[article.ID]; found {
			existingArticle.Tags = append(existingArticle.Tags, tags...)
		} else {
			articleMap[article.ID] = &article
		}
	}

	for _, article := range articleMap {
		articles = append(articles, *article)
	}

	return articles, nil
}

func (dao *articleDao) CountArticles(ctx context.Context) (int64, error) {
	var count int64
	result := dao.db.Model(model.Article{}).Count(&count) // 确保 Article 是正确的模型
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (dao *articleDao) ListArticleStatistics(ctx context.Context) ([]dto.ArticleStatisticsDTO, error) {
	var statistics []dto.ArticleStatisticsDTO
	query := `
	SELECT
		DATE_FORMAT(create_time, '%Y-%m-%d') AS date,
		COUNT(*) AS count
	FROM
		tb_article
	GROUP BY
		date
	ORDER BY
		date DESC`
	err := dao.db.Raw(query).Scan(&statistics).Error
	if err != nil {
		return nil, err
	}
	return statistics, nil
}

// GetArticleById 根据ID查询文章
func (dao *articleDao) GetArticleById(ctx context.Context, articleId int) (dto.ArticleDTO, error) {
	var article dto.ArticleDTO
	query := `
		SELECT
			a.id,
			a.article_cover,
			a.article_title,
			a.article_content,
			a.type,
			a.original_url,
			a.create_time,
			a.update_time,
			a.category_id,
			c.category_name,
			t.id AS tag_id,
			t.tag_name,
			a.user_id  -- 添加这行以查询文章的 UserID
		FROM tb_article a
		JOIN tb_category c ON a.category_id = c.id
		JOIN tb_article_tag atg ON a.id = atg.article_id
		JOIN tb_tag t ON t.id = atg.tag_id
		WHERE a.id = ?
		AND a.is_delete = 0
		AND a.status = 1`
	err := dao.db.Raw(query, articleId).Scan(&article).Error
	if err != nil {
		return dto.ArticleDTO{}, err
	}
	return article, nil
}

// ListArticleBacks 查询后台文章
func (dao *articleDao) ListArticleBacks(ctx context.Context, current, size int, condition vo.ConditionVO) ([]dto.ArticleBackDTO, error) {
	var articles []dto.ArticleBackDTO
	offset := (current - 1) * size
	query := `
		SELECT
			a.id,
			a.article_cover,
			a.article_title,
			a.type,
			a.is_top,
			a.is_delete,
			a.status,
			a.create_time,
			c.category_name,
			t.id AS tag_id,
			t.tag_name
		FROM (
			SELECT
				id,
				article_cover,
				article_title,
				type,
				is_top,
				is_delete,
				status,
				create_time,
				category_id
			FROM tb_article
			WHERE
				is_delete = ?
				<if test="condition.keywords != null">
					AND article_title LIKE ?
				</if>
				<if test="condition.status != null">
					AND status = ?
				</if>
				<if test="condition.categoryId != null">
					AND category_id = ?
				</if>
				<if test="condition.type != null">
					AND type = ?
				</if>
				<if test="condition.tagId != null">
					AND id IN (
						SELECT article_id FROM tb_article_tag WHERE tag_id = ?
					)
				</if>
			ORDER BY
				is_top DESC, id DESC
			LIMIT ?, ?
		) a
		LEFT JOIN tb_category c ON a.category_id = c.id
		LEFT JOIN tb_article_tag atg ON a.id = atg.article_id
		LEFT JOIN tb_tag t ON t.id = atg.tag_id
		ORDER BY
			is_top DESC, a.id DESC`

	// Build query parameters
	params := []interface{}{condition.IsDelete}
	if condition.Keywords != nil {
		params = append(params, fmt.Sprintf("%%%s%%", *condition.Keywords))
	}
	if condition.Status != nil {
		params = append(params, *condition.Status)
	}
	if condition.CategoryID != nil {
		params = append(params, *condition.CategoryID)
	}
	if condition.Type != nil {
		params = append(params, *condition.Type)
	}
	if condition.TagID != nil {
		params = append(params, *condition.TagID)
	}
	params = append(params, offset, size)

	err := dao.db.Raw(query, params...).Scan(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// CountArticleBacks 查询后台文章总量
func (dao *articleDao) CountArticleBacks(ctx context.Context, condition vo.ConditionVO) (int, error) {
	var count int64
	query := `
		SELECT COUNT(DISTINCT a.id)
		FROM tb_article a
		LEFT JOIN tb_article_tag tat ON a.id = tat.article_id
		WHERE
			is_delete = ?`

	// Build query parameters
	params := []interface{}{condition.IsDelete}

	// Append conditions dynamically
	if condition.Keywords != nil {
		query += " AND article_title LIKE ?"
		params = append(params, fmt.Sprintf("%%%s%%", *condition.Keywords))
	}
	if condition.Status != nil {
		query += " AND status = ?"
		params = append(params, *condition.Status)
	}
	if condition.CategoryID != nil {
		query += " AND category_id = ?"
		params = append(params, *condition.CategoryID)
	}
	if condition.Type != nil {
		query += " AND type = ?"
		params = append(params, *condition.Type)
	}
	if condition.TagID != nil {
		query += " AND tat.tag_id = ?"
		params = append(params, *condition.TagID)
	}

	err := dao.db.Raw(query, params...).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// ListRecommendArticles 查看推荐文章
func (dao *articleDao) ListRecommendArticles(ctx context.Context, articleId int) ([]dto.ArticleRecommendDTO, error) {
	var articles []dto.ArticleRecommendDTO
	query := `
		SELECT
			a.id,
			a.article_title,
			a.article_cover,
			a.create_time
		FROM (
			SELECT DISTINCT article_id
			FROM (
				SELECT tag_id FROM tb_article_tag WHERE article_id = ?
			) t
			JOIN tb_article_tag t1 ON t.tag_id = t1.tag_id
			WHERE article_id != ?
		) t2
		JOIN tb_article a ON t2.article_id = a.id
		WHERE a.is_delete = 0
		ORDER BY is_top DESC, a.id DESC
		LIMIT 6`

	err := dao.db.Raw(query, articleId, articleId).Scan(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (dao *articleDao) ListCategoryDTO(ctx context.Context) ([]dto.CategoryDTO, error) {
	var categories []dto.CategoryDTO
	query := `
		SELECT
			c.id, c.category_name, COUNT(a.id) AS article_count
		FROM
			tb_category c
		LEFT JOIN tb_article a ON c.id = a.category_id
		WHERE
			a.is_delete = 0
		GROUP BY
			c.id, c.category_name`
	err := dao.db.Raw(query).Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (dao *articleDao) ListTagDTO(ctx context.Context) ([]dto.TagDTO, error) {
	var tags []dto.TagDTO
	query := `
		SELECT
			t.id, t.tag_name
		FROM
			tb_tag t`
	err := dao.db.Raw(query).Scan(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetPreviousArticle 查询上一篇文章
func (dao *articleDao) GetPreviousArticle(ctx context.Context, articleId int) (model.Article, error) {
	var article model.Article
	result := dao.db.Table("tb_article").
		Select("id, article_title, article_cover").
		Where("is_delete = ?", 0).
		Where("status = ?", enums.PUBLIC.Status).
		Where("id < ?", articleId).
		Order("id DESC").
		Limit(1).
		Find(&article)

	if result.Error != nil {
		return model.Article{}, result.Error
	}

	return article, nil
}

// GetNextArticle 查询下一篇文章
func (dao *articleDao) GetNextArticle(ctx context.Context, articleId int) (model.Article, error) {
	var article model.Article
	result := dao.db.Table("tb_article").
		Select("id, article_title, article_cover").
		Where("is_delete = ?", 0).
		Where("status = ?", enums.PUBLIC.Status).
		Where("id > ?", articleId).
		Order("id ASC").
		Limit(1).
		Find(&article)

	if result.Error != nil {
		return model.Article{}, result.Error
	}

	return article, nil
}

func (dao *articleDao) ListArticlesByCondition(ctx context.Context, current int, size int, condition vo.ConditionVO) ([]dto.ArticlePreviewDTO, error) {
	log.Printf("ConditionVO before query: %+v", condition)

	var articles []dto.ArticlePreviewDTO

	// 确保 offset 不小于 0
	offset := (current - 1) * size
	if offset < 0 {
		offset = 0
	}

	query := `
    SELECT
        a.id,
        a.article_cover,
        a.article_title,
        a.create_time,
        a.category_id,
        c.category_name,
        t.id AS tag_id,
        t.tag_name
    FROM tb_article a
    JOIN tb_category c ON a.category_id = c.id
    JOIN tb_article_tag atg ON a.id = atg.article_id
    JOIN tb_tag t ON t.id = atg.tag_id
    WHERE
        a.is_delete = 0
        AND a.status = 1
        AND (? IS NULL OR a.category_id = ?)
        AND (? IS NULL OR t.id = ?)
    ORDER BY a.create_time DESC
    LIMIT ?, ?`

	params := []interface{}{
		condition.CategoryID, // 参数传递
		condition.CategoryID,
		condition.TagID, // 参数传递
		condition.TagID,
		offset, size,
	}

	// Debugging logs
	log.Println("About to execute query")
	log.Printf("Executing Query: %s", query)
	log.Printf("Parameters: %v", params)

	rows, err := dao.db.Raw(query, params...).Rows()
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	articleMap := make(map[int]*dto.ArticlePreviewDTO)

	for rows.Next() {
		var article dto.ArticlePreviewDTO
		var tag dto.TagDTO

		if err := rows.Scan(
			&article.ID,
			&article.ArticleCover,
			&article.ArticleTitle,
			&article.CreateTime,
			&article.CategoryId,
			&article.CategoryName,
			&tag.ID,
			&tag.TagName,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		if existingArticle, found := articleMap[article.ID]; found {
			existingArticle.TagDTOList = append(existingArticle.TagDTOList, tag)
		} else {
			article.TagDTOList = []dto.TagDTO{tag}
			articleMap[article.ID] = &article
		}
	}

	for _, article := range articleMap {
		articles = append(articles, *article)
	}

	log.Println("Query executed and results processed")
	return articles, nil
}

// SaveOrUpdate 保存或更新文章
func (dao *articleDao) SaveOrUpdate(ctx context.Context, article *model.Article) error {
	if article.ID == 0 {
		return dao.db.WithContext(ctx).Create(article).Error
	}
	return dao.db.WithContext(ctx).Save(article).Error
}

// 更新文章根据ID
func (dao *articleDao) UpdateById(ctx context.Context, article *model.Article) error {
	result := dao.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", article.ID).Updates(article)
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Updated rows: %d", result.RowsAffected) // 打印受影响的行数
	return nil
}

func (dao *articleDao) UpdateBatchById(ctx context.Context, articles []model.Article) error {
	// 遍历要更新的文章
	for _, article := range articles {
		// 更新文章状态，并手动设置 update_time 字段
		result := dao.db.Model(&model.Article{}).Where("id = ?", article.ID).Updates(map[string]interface{}{
			"is_delete":   article.IsDelete,
			"update_time": time.Now(),
		})

		// 检查是否更新成功
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// 删除文章
func (dao *articleDao) DeleteArticlesByIds(ctx context.Context, articleIdList []int) error {
	result := dao.db.WithContext(ctx).
		Where("id IN ?", articleIdList).
		Delete(&model.Article{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete articles: %v", result.Error)
	}
	return nil
}

// GetArticleCountByCategoryIDs 查询分类ID列表下的文章数量
func (dao *articleDao) GetArticleCountByCategoryIDs(ctx context.Context, categoryIDList []int) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("category_id IN ?", categoryIDList).
		Count(&count).Error
	return count, err
}

func (dao *articleDao) GetDb() *gorm.DB {
	return dao.db
}
