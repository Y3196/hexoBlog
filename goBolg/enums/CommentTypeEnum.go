package enums

type CommentTypeEnum struct {
	Type int
	Desc string
	Path string
}

var (
	ARTICLE = CommentTypeEnum{1, "文章评论", "/articles/"}
	LINK    = CommentTypeEnum{2, "友链评论", "/links/"}
	TALK    = CommentTypeEnum{3, "说说评论", "/talks/"}
)

var commentTypeEnums = []CommentTypeEnum{
	ARTICLE,
	LINK,
	TALK,
}

// GetCommentPath 获取评论路径
func GetCommentPath(commentType int) string {
	for _, v := range commentTypeEnums {
		if v.Type == commentType {
			return v.Path
		}
	}
	return ""
}

// GetCommentEnum 获取评论枚举
func GetCommentEnum(commentType int) *CommentTypeEnum {
	for _, v := range commentTypeEnums {
		if v.Type == commentType {
			return &v
		}
	}
	return nil
}
