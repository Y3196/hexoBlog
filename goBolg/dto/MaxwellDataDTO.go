package dto

// MaxwellDataDTO 代表 Maxwell 监听的数据 DTO
type MaxwellDataDTO struct {
	Database string                 `json:"database"` // 数据库
	Xid      int                    `json:"xid"`      // xid
	Data     map[string]interface{} `json:"data"`     // 数据
	Commit   bool                   `json:"commit"`   // 是否提交
	Type     string                 `json:"type"`     // 类型
	Table    string                 `json:"table"`    // 表
	Ts       int                    `json:"ts"`       // ts
}
