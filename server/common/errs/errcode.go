package errs

import "github.com/lpphub/golib/web"

/**
 * -1: 通用错误
 * 1000~1999: 系统错误
 * 2000~2999: 业务错误
 */
var (
	ErrServerError = web.Error{
		Code: -1,
		Msg:  "server internal error",
	}

	ErrApiFail = web.Error{
		Code: 1000,
		Msg:  "call api fail: %s",
	}

	ErrToast = web.Error{
		Code: 2000,
		Msg:  "%s",
	}
	ErrInvalidParam = web.Error{
		Code: 2001,
		Msg:  "invalid param",
	}
	ErrNotLogin = web.Error{
		Code: 2002,
		Msg:  "not login",
	}
	ErrInvalidToken = web.Error{
		Code: 2003,
		Msg:  "invalid token",
	}
)
