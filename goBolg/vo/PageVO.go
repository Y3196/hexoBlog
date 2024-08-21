package vo

type PageVO struct {
	ID        int    `json:"id" validate:"required" description:"页面id"`
	PageName  string `json:"pageName" validate:"required" description:"页面名称"`
	PageLabel string `json:"pageLabel" validate:"required" description:"页面标签"`
	PageCover string `json:"pageCover" validate:"required" description:"页面封面"`
}
