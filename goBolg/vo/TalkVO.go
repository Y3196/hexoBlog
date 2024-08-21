package vo

// TalkVO 代表说说对象
type TalkVO struct {
	ID      int    `json:"id"`                          // 说说id
	Content string `json:"content" validate:"required"` // 说说内容，不能为空
	Images  string `json:"images"`                      // 说说图片
	IsTop   int    `json:"isTop" validate:"required"`   // 置顶状态，不能为空
	Status  int    `json:"status" validate:"required"`  // 说说状态，不能为空
}
