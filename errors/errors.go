package errors

import "fmt"

const (
	UnknownCode   = -1
	UnknownRouter = -2
)

type Error struct {
	Code    int32
	Message string
	Reason  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s", e.Code, e.Reason, e.Message)
}

func New(code int, reason, message string) *Error {
	return &Error{
		Code:    int32(code),
		Message: message,
		Reason:  reason,
	}
}

func Newf(code int, reason, format string, a ...any) *Error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

func Errorf(code int, reason, format string, a ...any) error {
	return New(code, reason, fmt.Sprintf(format, a...))
}
