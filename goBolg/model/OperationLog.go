package model

import (
	"fmt"
	"time"
)

// OperationLog 代表操作日志的实体
type OperationLog struct {
	// 日志id
	ID int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`

	// 操作模块
	OptModule string `json:"optModule" gorm:"column:opt_module"`

	// 操作路径
	OptUrl string `json:"optUrl" gorm:"column:opt_url"`

	// 操作类型
	OptType string `json:"optType" gorm:"column:opt_type"`

	// 操作方法
	OptMethod string `json:"optMethod" gorm:"column:opt_method"`

	// 操作描述
	OptDesc string `json:"optDesc" gorm:"column:opt_desc"`

	// 请求方式
	RequestMethod string `json:"requestMethod" gorm:"column:request_method"`

	// 请求参数
	RequestParam string `json:"requestParam" gorm:"column:request_param"`

	// 返回数据
	ResponseData string `json:"responseData" gorm:"column:response_data"`

	// 用户id
	UserID int `json:"userId" gorm:"column:user_id"`

	// 用户昵称
	Nickname string `json:"nickname" gorm:"column:nickname"`

	// 用户登录ip
	IPAddress string `json:"ipAddress" gorm:"column:ip_address"`

	// ip来源
	IPSource string `json:"ipSource" gorm:"column:ip_source"`

	// 创建时间
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime;column:create_time"`

	// 修改时间
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime;column:update_time"`
}

// String 返回结构体的字符串表示
func (ol OperationLog) String() string {
	return fmt.Sprintf("OperationLog{ID=%d, OptModule='%s', OptUrl='%s', OptType='%s', OptMethod='%s', OptDesc='%s', RequestMethod='%s', RequestParam='%s', ResponseData='%s', UserID=%d, Nickname='%s', IPAddress='%s', IPSource='%s', CreateTime=%s, UpdateTime=%s}",
		ol.ID, ol.OptModule, ol.OptUrl, ol.OptType, ol.OptMethod, ol.OptDesc, ol.RequestMethod, ol.RequestParam, ol.ResponseData, ol.UserID, ol.Nickname, ol.IPAddress, ol.IPSource, ol.CreateTime.Format(time.RFC3339), ol.UpdateTime.Format(time.RFC3339))
}

// TableName 设置表名
func (OperationLog) TableName() string {
	return "tb_operation_log"
}
