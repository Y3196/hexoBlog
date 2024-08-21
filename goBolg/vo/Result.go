package vo

// Result 用于接口返回的结构体
type Result struct {
	Flag    bool        `json:"flag"`    // 返回状态，true 表示成功，false 表示失败
	Code    int         `json:"code"`    // 返回码
	Message string      `json:"message"` // 返回信息
	Data    interface{} `json:"data"`    // 返回数据，支持任意类型
}

// Ok 创建一个表示成功的结果，没有数据
func Ok() Result {
	return Result{
		Flag:    true,
		Code:    200,
		Message: "操作成功",
		Data:    nil,
	}
}

// OkWithData 创建一个表示成功的结果，包含数据
func OkWithData(data interface{}) Result {
	return Result{
		Flag:    true,
		Code:    200,
		Message: "操作成功",
		Data:    data,
	}
}

// OkWithMessage 创建一个表示成功的结果，包含自定义消息
func OkWithMessage(message string) Result {
	return Result{
		Flag:    true,
		Code:    200,
		Message: message,
		Data:    nil,
	}
}

// OkWithDataAndMessage 创建一个表示成功的结果，包含数据和自定义消息
func OkWithDataAndMessage(data interface{}, message string) Result {
	return Result{
		Flag:    true,
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// Fail 创建一个表示失败的结果，没有数据
func Fail() Result {
	return Result{
		Flag:    false,
		Code:    500,
		Message: "操作失败",
		Data:    nil,
	}
}

// FailWithData 创建一个表示失败的结果，包含数据
func FailWithData(data interface{}) Result {
	return Result{
		Flag:    false,
		Code:    500,
		Message: "操作失败",
		Data:    data,
	}
}

// FailWithMessage 创建一个表示失败的结果，包含自定义消息
func FailWithMessage(message string) Result {
	return Result{
		Flag:    false,
		Code:    500,
		Message: message,
		Data:    nil,
	}
}

// FailWithDataAndMessage 创建一个表示失败的结果，包含数据和自定义消息
func FailWithDataAndMessage(data interface{}, message string) Result {
	return Result{
		Flag:    false,
		Code:    500,
		Message: message,
		Data:    data,
	}
}

// FailWithCodeAndMessage 创建一个表示失败的结果，包含自定义码和消息
func FailWithCodeAndMessage(code int, message string) Result {
	return Result{
		Flag:    false,
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
