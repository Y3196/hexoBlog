package enums

// ArticleStatusEnum 文章状态枚举
type ArticleStatusEnum struct {
	Status int
	Desc   string
}

// 定义文章状态常量
var (
	PUBLIC = ArticleStatusEnum{1, "公开"}
	SECRET = ArticleStatusEnum{2, "私密"}
	DRAFT  = ArticleStatusEnum{3, "草稿"}
)

// GetArticleStatusEnums 获取所有文章状态枚举
func GetArticleStatusEnums() []ArticleStatusEnum {
	return []ArticleStatusEnum{PUBLIC, SECRET, DRAFT}
}
