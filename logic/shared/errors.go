package shared

var (
	ErrServerError = NewError(500, "server internal error")

	ErrNoToken      = NewError(1000, "no token")
	ErrInvalidToken = NewError(1001, "invalid token")

	ErrInvalidParam   = NewError(1100, "参数错误")
	ErrRecordNotFound = NewError(1101, "数据不存在")

	ErrUserExists       = NewError(2101, "用户已存在")
	ErrEmailExists      = NewError(2102, "邮箱已存在")
	ErrInvalidPassword  = NewError(2103, "密码格式错误")
	ErrRoleExists       = NewError(2110, "角色已存在")
	ErrPermissionExists = NewError(2120, "权限已存在")
)

type Error struct {
	Code int
	Msg  string
}

func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Msg:  msg,
	}
}

func (err Error) Error() string {
	return err.Msg
}
