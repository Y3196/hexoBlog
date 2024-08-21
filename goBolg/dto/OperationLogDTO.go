package dto

import "time"

// OperationLogDTO 代表操作日志 DTO
type OperationLogDTO struct {
	ID            int       `json:"id"`            // 日志id
	OptModule     string    `json:"optModule"`     // 操作模块
	OptUrl        string    `json:"optUrl"`        // 操作路径
	OptType       string    `json:"optType"`       // 操作类型
	OptMethod     string    `json:"optMethod"`     // 操作方法
	OptDesc       string    `json:"optDesc"`       // 操作描述
	RequestMethod string    `json:"requestMethod"` // 请求方式
	RequestParam  string    `json:"requestParam"`  // 请求参数
	ResponseData  string    `json:"responseData"`  // 返回数据
	Nickname      string    `json:"nickname"`      // 用户昵称
	IpAddress     string    `json:"ipAddress"`     // 用户登录ip
	IpSource      string    `json:"ipSource"`      // ip来源
	CreateTime    time.Time `json:"createTime"`    // 创建时间
}
