package enums

// FilePathEnum represents the file path enumeration
type FilePathEnum struct {
	Path string
	Desc string
}

// Predefined file path enums
var (
	Avatar   = FilePathEnum{"avatar/", "头像路径"}
	Article  = FilePathEnum{"articles/", "文章图片路径"}
	Voice    = FilePathEnum{"voice/", "音频路径"}
	Photo    = FilePathEnum{"photos/", "相册路径"}
	Config   = FilePathEnum{"config/", "配置图片路径"}
	Talk     = FilePathEnum{"talks/", "说说图片路径"}
	Markdown = FilePathEnum{"markdown/", "md文件路径"}
)

// GetFilePathEnums returns a list of all file path enums
func GetFilePathEnums() []FilePathEnum {
	return []FilePathEnum{
		Avatar,
		Article,
		Voice,
		Photo,
		Config,
		Talk,
		Markdown,
	}
}
