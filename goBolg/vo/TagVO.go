package vo

// TagVO 代表标签对象
type TagVO struct {
	ID      int    `json:"id"`                          // 标签id
	TagName string `json:"tagName" validate:"required"` // 标签名
}
