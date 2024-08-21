package dao

import (
	"context"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
)

// TalkDao 说说 DAO 接口
type TalkDao interface {
	// 获取说说列表
	ListTalks(ctx context.Context, current, size int) ([]dto.TalkDTO, error)

	// 查看后台说说
	ListBackTalks(ctx context.Context, current, size int, condition vo.ConditionVO) ([]dto.TalkBackDTO, error)

	// 根据id查看说说
	GetTalkById(ctx context.Context, talkId int) (*dto.TalkDTO, error)

	// 根据id查看后台说说
	GetBackTalkById(ctx context.Context, talkId uint) (*dto.TalkBackDTO, error)

	// 查询最新10条说说
	ListLatestTalks(ctx context.Context) ([]string, error)

	// 查询说说总量
	CountTalks(ctx context.Context, status *int) (int64, error)

	// 删除说说
	DeleteTalks(ctx context.Context, talkIdList []uint) error

	CountAll(ctx context.Context) (int, error)
}

type talkDao struct {
	db *gorm.DB
}

// NewTalkDao 创建新的 TalkDao 实例
func NewTalkDao(db *gorm.DB) TalkDao {
	return &talkDao{db: db}
}

// ListTalks 获取说说列表
func (dao *talkDao) ListTalks(ctx context.Context, current, size int) ([]dto.TalkDTO, error) {
	var talks []dto.TalkDTO
	err := dao.db.WithContext(ctx).
		Raw(`
			SELECT
				t.id,
				nickname,
				avatar,
				content,
				images,
				t.is_top,
				t.create_time
			FROM
				tb_talk t
				JOIN tb_user_info ui ON t.user_id = ui.id
			WHERE
				t.status = 1
			ORDER BY
				t.is_top DESC,
				t.id DESC
			LIMIT ?, ?
		`, current, size).
		Scan(&talks).Error

	if err != nil {
		return nil, err
	}
	return talks, nil
}

// ListBackTalks 查看后台说说
func (dao *talkDao) ListBackTalks(ctx context.Context, current, size int, condition vo.ConditionVO) ([]dto.TalkBackDTO, error) {
	var talks []dto.TalkBackDTO
	query := dao.db.WithContext(ctx).
		Raw(`
			SELECT
				t.id,
				nickname,
				avatar,
				content,
				images,
				t.is_top,
				t.status,
				t.create_time
			FROM
				tb_talk t
				JOIN tb_user_info ui ON t.user_id = ui.id
			WHERE
				1 = 1
		`)

	if condition.Status != nil {
		query = query.Where("t.status = ?", *condition.Status)
	}

	query = query.Order("t.is_top DESC, t.id DESC").
		Offset(current).
		Limit(size).
		Find(&talks)

	if query.Error != nil {
		return nil, query.Error
	}
	return talks, nil
}

// GetTalkById 根据id查看说说
func (dao *talkDao) GetTalkById(ctx context.Context, talkId int) (*dto.TalkDTO, error) {
	var talk dto.TalkDTO
	err := dao.db.WithContext(ctx).
		Raw(`
			SELECT
				t.id,
				nickname,
				avatar,
				content,
				images,
				t.create_time
			FROM
				tb_talk t
				JOIN tb_user_info ui ON t.user_id = ui.id
			WHERE
				t.id = ?
				AND t.status = 1
		`, talkId).
		Scan(&talk).Error

	if err != nil {
		return nil, err
	}
	return &talk, nil
}

// GetBackTalkById 根据id查看后台说说
func (dao *talkDao) GetBackTalkById(ctx context.Context, talkId uint) (*dto.TalkBackDTO, error) {
	var talk dto.TalkBackDTO
	err := dao.db.WithContext(ctx).
		Raw(`
			SELECT
				t.id,
				nickname,
				avatar,
				content,
				images,
				t.is_top,
				t.status,
				t.create_time
			FROM
				tb_talk t
				JOIN tb_user_info ui ON t.user_id = ui.id
			WHERE
				t.id = ?
		`, talkId).
		Scan(&talk).Error

	if err != nil {
		return nil, err
	}
	return &talk, nil
}

// ListLatestTalks 查询最新10条说说
func (dao *talkDao) ListLatestTalks(ctx context.Context) ([]string, error) {
	var talks []model.Talk
	err := dao.db.WithContext(ctx).
		Where("status = ?", 1).
		Order("is_top DESC, id DESC").
		Limit(10).
		Find(&talks).Error

	if err != nil {
		return nil, err
	}

	var contents []string
	for _, talk := range talks {
		content := talk.Content
		if len(content) > 200 {
			content = content[:200] // 取前200个字符
		}
		// 假设 HTMLUtils.DeleteHTMLTag 处理 HTML 标签
		contents = append(contents, content) // 这里省略了删除 HTML 标签的处理
	}
	return contents, nil
}

// CountTalks 查询说说总量
func (dao *talkDao) CountTalks(ctx context.Context, status *int) (int64, error) {
	var count int64
	query := dao.db.WithContext(ctx).Model(&model.Talk{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteTalks 删除说说
func (dao *talkDao) DeleteTalks(ctx context.Context, talkIdList []uint) error {
	return dao.db.WithContext(ctx).Where("id IN ?", talkIdList).Delete(&model.Talk{}).Error
}

func (dao *talkDao) CountAll(ctx context.Context) (int, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.Talk{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
