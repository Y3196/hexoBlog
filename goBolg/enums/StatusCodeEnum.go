package enums

type StatusCodeEnum struct {
	Code int
	Desc string
}

var (
	SUCCESS            = StatusCodeEnum{20000, "操作成功"}
	NO_LOGIN           = StatusCodeEnum{40001, "用户未登录"}
	AUTHORIZED         = StatusCodeEnum{40300, "没有操作权限"}
	SYSTEM_ERROR       = StatusCodeEnum{50000, "系统异常"}
	FAIL               = StatusCodeEnum{51000, "操作失败"}
	VALID_ERROR        = StatusCodeEnum{52000, "参数格式不正确"}
	USERNAME_EXIST     = StatusCodeEnum{52001, "用户名已存在"}
	USERNAME_NOT_EXIST = StatusCodeEnum{52002, "用户名不存在"}
	QQ_LOGIN_ERROR     = StatusCodeEnum{53001, "qq登录错误"}
	WEIBO_LOGIN_ERROR  = StatusCodeEnum{53002, "微博登录错误"}
)

// StatusCodeEnums 是一个所有状态码枚举值的列表，可以用来进行查找和迭代
var StatusCodeEnums = []StatusCodeEnum{
	SUCCESS,
	NO_LOGIN,
	AUTHORIZED,
	SYSTEM_ERROR,
	FAIL,
	VALID_ERROR,
	USERNAME_EXIST,
	USERNAME_NOT_EXIST,
	QQ_LOGIN_ERROR,
	WEIBO_LOGIN_ERROR,
}
