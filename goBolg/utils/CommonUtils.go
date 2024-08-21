package utils

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

// CheckEmail 检测邮箱是否合法
func CheckEmail(email string) bool {
	rule := `^\w+((-|\.)\w+)*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$`
	re := regexp.MustCompile(rule)
	return re.MatchString(email)
}

// GetBracketsContent 获取括号内容
func GetBracketsContent(str string) string {
	start := strings.Index(str, "(")
	end := strings.Index(str, ")")
	if start != -1 && end != -1 && end > start {
		return str[start+1 : end]
	}
	return ""
}

// GetRandomCode 生成6位随机验证码
func GetRandomCode() string {
	code := make([]byte, 6)
	for i := 0; i < 6; i++ {
		code[i] = byte(rand.Intn(10) + '0')
	}
	return string(code)
}

// CastList 转换为列表
func CastList[T any](obj interface{}, constructor func(interface{}) (T, error)) ([]T, error) {
	var result []T
	if list, ok := obj.([]interface{}); ok {
		for _, item := range list {
			elem, err := constructor(item)
			if err != nil {
				return nil, err
			}
			result = append(result, elem)
		}
	}
	return result, nil
}

// CastSet 转换为集合
func CastSet[T comparable](obj interface{}, constructor func(interface{}) (T, error)) (map[T]struct{}, error) {
	result := make(map[T]struct{})
	if set, ok := obj.([]interface{}); ok {
		for _, item := range set {
			elem, err := constructor(item)
			if err != nil {
				return nil, err
			}
			result[elem] = struct{}{}
		}
	}
	return result, nil
}

// 示例构造函数，用于将接口转换为目标类型
func IntConstructor(item interface{}) (int, error) {
	if v, ok := item.(float64); ok {
		return int(v), nil
	}
	return 0, errors.New("type assertion to int failed")
}

// StringConstructor 用于将接口转换为字符串
func StringConstructor(item interface{}) (string, error) {
	if v, ok := item.(string); ok {
		return v, nil
	}
	return "", errors.New("type assertion to string failed")
}
