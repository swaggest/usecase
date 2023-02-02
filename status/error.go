package status

import (
	"errors"
	"strings"
)

var codeToMsg = func() map[Code]string {
	res := make(map[Code]string, len(strToCode))
	for str, code := range strToCode {
		res[code] = strings.ToLower(strings.ReplaceAll(str, "_", " "))
	}

	return res
}()

// Error returns string value of status code.
func (c Code) Error() string {
	return codeToMsg[c]
}

type errorWithStatus struct {
	err  error
	code Code
}

func (e errorWithStatus) Error() string {
	return codeToMsg[e.code] + ": " + e.err.Error()
}

func (e errorWithStatus) Unwrap() error {
	return e.err
}

func (e errorWithStatus) Is(target error) bool {
	if target == e.code { //nolint:goerr113 // Target is expected to be plain status error.
		return true
	}

	return errors.Is(e.err, target)
}

func (e errorWithStatus) Status() Code {
	return e.code
}

// Wrap adds canonical status to error.
func Wrap(err error, code Code) error {
	if err == nil {
		return code
	}

	return errorWithStatus{
		err:  err,
		code: code,
	}
}
