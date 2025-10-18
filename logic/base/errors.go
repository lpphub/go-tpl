package base

var (
	ErrServerError = NewError(-1, "server internal error")

	ErrNoToken      = NewError(1000, "no token")
	ErrInvalidToken = NewError(1002, "invalid token")
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
