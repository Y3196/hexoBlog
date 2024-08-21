package utils

import (
	"github.com/importcjj/sensitive"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
)

// 初始化敏感词过滤器
var filter = sensitive.New()

func init() {
	filter.LoadWordDict("path/to/sensitive-words.txt") // 设置敏感词字典路径
}

// HTMLUtils 处理HTML标签和敏感词
type HTMLUtils struct{}

// Filter 删除标签和敏感词过滤
func (HTMLUtils) Filter(source string) string {
	// 敏感词过滤
	source = filter.Replace(source, '*')

	// 保留图片标签
	re := regexp.MustCompile(`<[^>]*>`)
	source = re.ReplaceAllStringFunc(source, func(tag string) string {
		if match, _ := regexp.MatchString(`^<img`, tag); match {
			return tag
		}
		return ""
	})

	re = regexp.MustCompile(`onload\s*=\s*['"].*?['"]`)
	source = re.ReplaceAllString(source, "")

	re = regexp.MustCompile(`onerror\s*=\s*['"].*?['"]`)
	source = re.ReplaceAllString(source, "")

	return deleteHTMLTag(source)
}

// deleteHTMLTag 删除HTML标签
func deleteHTMLTag(source string) string {
	// 删除转义字符
	re := regexp.MustCompile(`&.{2,6}?;`)
	source = re.ReplaceAllString(source, "")

	// 删除script标签
	re = regexp.MustCompile(`<[\\s]*?script[^>]*?>[\\s\\S]*?<[\\s]*?\\/[\\s]*?script[\\s]*?>`)
	source = re.ReplaceAllString(source, "")

	// 删除style标签
	re = regexp.MustCompile(`<[\\s]*?style[^>]*?>[\\s\\S]*?<[\\s]*?\\/[\\s]*?style[\\s]*?>`)
	source = re.ReplaceAllString(source, "")

	// 使用 bluemonday 进一步清理 HTML
	p := bluemonday.UGCPolicy()
	return p.Sanitize(source)
}

// HTMLFilter 是对 HTMLUtils 的 Filter 方法的简单包装
func HTMLFilter(source string) string {
	return HTMLUtils{}.Filter(source)
}
