package dao

import (
	"context"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
	"log"
	"time"
)

// FriendLinkDao 友情链接 DAO 接口
type FriendLinkDao interface {
	// 获取所有友情链接列表
	AllFriendLinks(ctx context.Context) ([]model.FriendLink, error)

	// 分页查询友情链接
	PagedFriendLinks(ctx context.Context, page int, size int, condition vo.ConditionVO) ([]model.FriendLink, int64, error)

	SelectList(ctx context.Context) ([]model.FriendLink, error)

	Save(ctx context.Context, friendLink *model.FriendLink) error
	Update(ctx context.Context, friendLink *model.FriendLink) error

	FindById(ctx context.Context, id uint) (*model.FriendLink, error)

	RemoveByIds(ctx context.Context, ids []uint) error
}

type friendLinkDao struct {
	db *gorm.DB
}

// NewFriendLinkDao 创建新的 FriendLinkDao 实例
func NewFriendLinkDao(db *gorm.DB) FriendLinkDao {
	return &friendLinkDao{db: db}
}

// AllFriendLinks 获取所有友情链接列表
func (dao *friendLinkDao) AllFriendLinks(ctx context.Context) ([]model.FriendLink, error) {
	var friendLinks []model.FriendLink

	err := dao.db.WithContext(ctx).Find(&friendLinks).Error
	if err != nil {
		return nil, err
	}

	return friendLinks, nil
}

// PagedFriendLinks 分页查询友情链接
func (dao *friendLinkDao) PagedFriendLinks(ctx context.Context, page int, size int, condition vo.ConditionVO) ([]model.FriendLink, int64, error) {
	var friendLinks []model.FriendLink
	var count int64

	offset := (page - 1) * size
	query := dao.db.WithContext(ctx).Model(&model.FriendLink{}).Offset(offset).Limit(size)

	if condition.Keywords != nil && *condition.Keywords != "" {
		query = query.Where("link_name LIKE ?", "%"+*condition.Keywords+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Find(&friendLinks).Error
	if err != nil {
		return nil, 0, err
	}

	return friendLinks, count, nil
}

func (dao *friendLinkDao) SelectList(ctx context.Context) ([]model.FriendLink, error) {
	var friendLinks []model.FriendLink
	err := dao.db.WithContext(ctx).Unscoped().Find(&friendLinks).Error
	if err != nil {
		return nil, err
	}
	return friendLinks, nil
}

func (dao *friendLinkDao) Save(ctx context.Context, friendLink *model.FriendLink) error {
	log.Println("Creating new friend link")
	result := dao.db.WithContext(ctx).Omit("update_time").Create(friendLink)
	log.Printf("SQL: %v", result.Statement.SQL.String())
	return result.Error
}

func (dao *friendLinkDao) Update(ctx context.Context, friendLink *model.FriendLink) error {
	log.Println("Updating existing friend link")
	now := time.Now()
	friendLink.UpdateTime = &now
	result := dao.db.WithContext(ctx).Save(friendLink)
	log.Printf("SQL: %v", result.Statement.SQL.String())
	return result.Error
}

func (dao *friendLinkDao) FindById(ctx context.Context, id uint) (*model.FriendLink, error) {
	var friendLink model.FriendLink
	err := dao.db.WithContext(ctx).First(&friendLink, id).Error
	return &friendLink, err
}

func (dao *friendLinkDao) RemoveByIds(ctx context.Context, ids []uint) error {
	return dao.db.WithContext(ctx).Delete(&model.FriendLink{}, ids).Error
}
