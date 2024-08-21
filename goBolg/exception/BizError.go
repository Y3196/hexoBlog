// rabbitService/error.go
package exception

import (
	"fmt"
	"goBolg/enums"
)

type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewBizError 创建一个新的业务错误实例
func NewBizError(code int, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

func NewBizException(message string) *BizError {
	return &BizError{
		Code:    enums.FAIL.Code,
		Message: message,
	}
}
