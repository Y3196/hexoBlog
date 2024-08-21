package dao

import (
	"context"
	"goBolg/model"
	"gorm.io/gorm"
)

type WebsiteConfigDao interface {
	SelectByID(ctx context.Context, db *gorm.DB, id uint) (*model.WebsiteConfig, error)
	UpdateByID(ctx context.Context, db *gorm.DB, config *model.WebsiteConfig) error
}

type websiteConfigDao struct {
	db *gorm.DB
}

func NewWebsiteConfigDao(db *gorm.DB) WebsiteConfigDao {
	return &websiteConfigDao{db: db}
}

func (dao *websiteConfigDao) SelectByID(ctx context.Context, db *gorm.DB, id uint) (*model.WebsiteConfig, error) {
	var config model.WebsiteConfig
	result := db.First(&config, id)
	return &config, result.Error
}

// UpdateByID 更新数据库中的网站配置
func (dao *websiteConfigDao) UpdateByID(ctx context.Context, db *gorm.DB, config *model.WebsiteConfig) error {
	result := db.Model(&model.WebsiteConfig{}).Where("id = ?", config.ID).Updates(config)
	return result.Error
}
