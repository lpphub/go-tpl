package errs

import "github.com/lpphub/golib/web"

var (
	ErrServerError = web.Error{
		Code: -1,
		Msg:  "server internal error",
	}
	ErrNotLogin = web.Error{
		Code: 1001,
		Msg:  "not login",
	}
	ErrInvalidToken = web.Error{
		Code: 1002,
		Msg:  "invalid token",
	}
	ErrParamInvalid = web.Error{
		Code: 1100,
		Msg:  "invalid param",
	}
	ErrApiFail = web.Error{
		Code: 1101,
		Msg:  "call api fail: %s",
	}
	ErrToast = web.Error{
		Code: 1199,
		Msg:  "%s",
	}
)
