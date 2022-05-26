package usecase

import "github.com/swaggest/usecase/status"

// Error is an error with contextual information.
type Error struct {
	AppCode    int
	StatusCode status.Code
	Value      error
	Context    map[string]interface{}
}

// Error returns error message.
func (e Error) Error() string {
	return e.Unwrap().Error()
}

// Fields exposes structured context of error.
func (e Error) Fields() map[string]interface{} {
	return e.Context
}

// AppErrCode returns application level error code.
func (e Error) AppErrCode() int {
	return e.AppCode
}

// Status returns status code of error.
func (e Error) Status() status.Code {
	return e.StatusCode
}

// Unwrap returns parent error.
func (e Error) Unwrap() error {
	if e.StatusCode != 0 {
		return status.Wrap(e.Value, e.StatusCode)
	}

	return e.Value
}

type sentinelError string

func (s sentinelError) Error() string {
	return string(s)
}
