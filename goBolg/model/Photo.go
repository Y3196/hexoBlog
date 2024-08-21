package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// Photo 代表照片的实体
type Photo struct {
	// 照片id
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 相册id
	AlbumID int `json:"albumId" gorm:"column:album_id"`

	// 照片名
	PhotoName string `json:"photoName" gorm:"column:photo_name;type:varchar(100)"`

	// 照片描述
	PhotoDesc string `json:"photoDesc" gorm:"column:photo_desc"`

	// 照片地址
	PhotoSrc string `json:"photoSrc" gorm:"column:photo_src"`

	// 是否删除
	IsDelete int `json:"isDelete" gorm:"column:is_delete"`

	// 创建时间
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`

	// 修改时间
	UpdateTime *time.Time `gorm:"column:update_time"`
}

// String 返回结构体的字符串表示
func (p Photo) String() string {
	return fmt.Sprintf("Photo{ID=%d, AlbumID=%d, PhotoName='%s', PhotoDesc='%s', PhotoSrc='%s', IsDelete=%d, CreateTime=%s, UpdateTime=%s}",
		p.ID, p.AlbumID, p.PhotoName, p.PhotoDesc, p.PhotoSrc, p.IsDelete, p.CreateTime.Format(time.RFC3339), p.UpdateTime.Format(time.RFC3339))
}

func (Photo) TableName() string {
	return "tb_photo"
}

// BeforeCreate 在创建记录之前设置创建时间
func (a *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	a.CreateTime = time.Now()
	a.UpdateTime = nil // 确保在创建时不设置更新时间
	return
}

// BeforeUpdate 在更新记录之前设置更新时间
func (a *Photo) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	a.UpdateTime = &now
	return
}
