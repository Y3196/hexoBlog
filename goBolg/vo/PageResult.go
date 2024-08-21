package vo

// PageResult 分页结果结构体
type PageResult[T any] struct {
	RecordList []T `json:"recordList"` // 分页列表
	Count      int `json:"count"`      // 总数
	//Records    []dto.CommentDTO
}

// NewPageResult 构造函数，创建一个新的分页结果
func NewPageResult[T any](recordList []T, count int) PageResult[T] {
	return PageResult[T]{
		RecordList: recordList,
		Count:      count,
	}
}
