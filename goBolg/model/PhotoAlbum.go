package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// PhotoAlbum 代表相册的实体
type PhotoAlbum struct {
	// 主键
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 相册名
	AlbumName string `json:"albumName" gorm:"column:album_name"`

	// 相册描述
	AlbumDesc string `json:"albumDesc" gorm:"column:album_desc"`

	// 相册封面
	AlbumCover string `json:"albumCover" gorm:"column:album_cover"`

	// 是否删除
	IsDelete int `json:"isDelete" gorm:"column:is_delete"`

	// 状态值 1公开 2私密
	Status int `json:"status" gorm:"column:status"`

	// 创建时间
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime;column:create_time"`

	// 修改时间
	UpdateTime *time.Time `gorm:"column:update_time"`
}

// BeforeCreate 在创建记录之前设置创建时间
func (pa *PhotoAlbum) BeforeCreate(tx *gorm.DB) (err error) {
	if pa.CreateTime.IsZero() {
		pa.CreateTime = time.Now()
	}
	return
}

// String 返回结构体的字符串表示
func (pa PhotoAlbum) String() string {
	return fmt.Sprintf("PhotoAlbum{ID=%d, AlbumName='%s', AlbumDesc='%s', AlbumCover='%s', IsDelete=%d, Status=%d, CreateTime=%s, UpdateTime=%s}",
		pa.ID, pa.AlbumName, pa.AlbumDesc, pa.AlbumCover, pa.IsDelete, pa.Status, pa.CreateTime.Format(time.RFC3339), pa.UpdateTime.Format(time.RFC3339))
}

func (PhotoAlbum) TableName() string {
	return "tb_photo_album"
}
