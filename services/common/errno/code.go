package errno

// 业务错误码规范：5位
// 服务级别错误码：1 位数进行表示，比如 1 为系统级错误；2 为普通错误，通常是由用户非法操作引起。
// 模块级错误码：2 位数进行表示，比如 01 为用户模块；02 为订单模块。
// 具体错误码：2 位数进行表示，比如 01 为手机号不合法；02 为验证码输入错误。

// 其中
// 0：成功
// -1：未知异常

var (
	OK = NewError(0, "OK")

	// 服务级错误码
	ErrServer    = NewError(10001, "服务异常，请联系管理员")
	ErrParam     = NewError(10002, "参数有误")
	ErrSignParam = NewError(10003, "签名参数有误")

	// 模块级错误码 - 用户模块
	ErrUserPhone          = NewError(20101, "用户手机号不合法")
	ErrUserCaptcha        = NewError(20102, "用户验证码有误")
	ErrUserPassword       = NewError(20103, "用户密码有误")
	ErrUserNameNotUnique  = NewError(20104, "用户名不唯一")
	ErrUserIDNotExit      = NewError(20105, "用户不存在")
	ErrUserNameOrPassword = NewError(20106, "用户名或密码有误")
)
