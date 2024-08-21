package utils

import "encoding/json"

// DecodeDoubleEscapedJSON 处理双重转义的JSON字符串
func DecodeDoubleEscapedJSON(data string) (string, error) {
	var intermediateStr string
	// 第一次解码
	err := json.Unmarshal([]byte(data), &intermediateStr)
	if err != nil {
		return "", err
	}
	// 将解码后的字符串再次解码
	return intermediateStr, nil
}
