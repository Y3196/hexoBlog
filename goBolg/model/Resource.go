package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// Resource 代表资源的实体
type Resource struct {
	// 权限id
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 权限名
	ResourceName string `json:"resourceName" gorm:"column:resource_name"`

	// 权限路径
	URL string `json:"url" gorm:"column:url"`

	// 请求方式
	RequestMethod string `json:"requestMethod" gorm:"column:request_method"`

	// 父权限id
	ParentID *uint `json:"parentId" gorm:"column:parent_id"`

	// 是否匿名访问
	IsAnonymous int `json:"isAnonymous" gorm:"column:is_anonymous"`

	// 创建时间
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime;column:create_time"`

	// 修改时间
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time"`
}

// BeforeCreate 在创建记录之前设置创建时间
func (r *Resource) BeforeCreate(tx *gorm.DB) (err error) {
	if r.CreateTime.IsZero() {
		r.CreateTime = time.Now()
	}
	return
}

// String 返回结构体的字符串表示
func (r Resource) String() string {
	return fmt.Sprintf("Resource{ID=%d, ResourceName='%s', URL='%s', RequestMethod='%s', ParentID=%d, IsAnonymous=%d, CreateTime=%s, UpdateTime=%s}",
		r.ID, r.ResourceName, r.URL, r.RequestMethod, r.ParentID, r.IsAnonymous, r.CreateTime.Format(time.RFC3339), r.UpdateTime.Format(time.RFC3339))
}

func (Resource) TableName() string {
	return "tb_resource"
}
