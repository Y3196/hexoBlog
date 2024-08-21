package vo

import "time"

type ConditionVO struct {
	Current    *int       `form:"current"`    // 页码
	Size       *int       `form:"size"`       // 页大小
	Keywords   *string    `form:"keywords"`   // 关键词
	CategoryID *int       `form:"categoryId"` // 分类ID
	TagID      *int       `form:"tagId"`      // 标签ID
	AlbumID    *int       `form:"albumId"`    // 专辑ID
	LoginType  *string    `form:"loginType"`  // 登录类型
	Type       *int       `form:"type"`       // 类型
	Status     *int       `form:"status"`     // 状态
	StartTime  *time.Time `form:"startTime"`  // 开始时间
	EndTime    *time.Time `form:"endTime"`    // 结束时间
	IsDelete   *bool      `form:"isDelete"`   // 是否删除
	IsReview   *int       `form:"isReview"`   // 是否审核
}
