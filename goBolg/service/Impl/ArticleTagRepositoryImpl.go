package Impl

import (
	"gorm.io/gorm"
)

// ArticleTagRepositoryImpl 是 ArticleTagRepository 接口的实现
type ArticleTagRepositoryImpl struct {
	db *gorm.DB
}

// NewArticleTagRepositoryImpl 创建一个新的 ArticleTagRepositoryImpl 实例
func NewArticleTagRepositoryImpl(db *gorm.DB) *ArticleTagRepositoryImpl {
	return &ArticleTagRepositoryImpl{db: db}
}
