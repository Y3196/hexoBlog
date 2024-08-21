package constants

type ContextKey string

const (
	// 否
	False = 0

	// 是
	True = 1

	// 高亮标签
	PreTag  = "<span style='color:#f47466'>"
	PostTag = "</span>"

	// 当前页码
	Current = "current"

	// 页码条数
	Size = "size"

	// 博主id
	BloggerID = 1

	// 默认条数
	DefaultSize = "10"

	// 默认用户昵称
	DefaultNickname = "用户"

	// 浏览文章集合
	ArticleSet = "articleSet"

	// 前端组件名
	Component = "Layout"

	// 省
	Province = "省"

	// 市
	City = "市"

	// 未知的
	Unknown = "未知"

	// JSON 格式
	ApplicationJSON = "application/json;charset=utf-8"

	// 默认的配置id
	DefaultConfigID = 1

	//
	UserContextKey ContextKey = "user"
)
