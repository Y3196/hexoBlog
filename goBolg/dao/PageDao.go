package dao

import (
	"context"
	"goBolg/model"
	"gorm.io/gorm"
)

type PageDao interface {
	SaveOrUpdate(ctx context.Context, db *gorm.DB, page *model.Page) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*model.Page, error)
}

type pageDao struct {
	db *gorm.DB
}

func NewPageDao(db *gorm.DB) PageDao {
	return &pageDao{db: db}
}

func (dao *pageDao) SaveOrUpdate(ctx context.Context, db *gorm.DB, page *model.Page) error {
	return db.Save(page).Error
}

func (dao *pageDao) Delete(ctx context.Context, id uint) error {
	return dao.db.Delete(&model.Page{}, id).Error
}

func (dao *pageDao) List(ctx context.Context) ([]*model.Page, error) {
	var pages []*model.Page
	result := dao.db.Find(&pages)
	return pages, result.Error
}
