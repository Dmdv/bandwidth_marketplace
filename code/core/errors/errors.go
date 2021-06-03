package errors

import (
	"os"
)

const (
	delim = ": "
)

type (
	// Error type for a new application error.
	Error struct {
		Code string `json:"code,omitempty"`
		Msg  string `json:"msg"`
	}
)

// Error implements error interface.
func (err *Error) Error() string {
	return err.Code + delim + err.Msg
}

// NewError creates a new error.
func NewError(code string, msg string) *Error {
	return &Error{Code: code, Msg: msg}
}

// ExitErr prints error to os.Stderr and call os.Exit with given code.
func ExitErr(msg string, err error, code int) {
	msg = WrapError("exit", msg, err).Error()
	_, _ = os.Stderr.Write([]byte(msg))
	os.Exit(code)
}

// ExitMsg prints message to os.Stderr and call os.Exit with given code.
func ExitMsg(msg string, code int) {
	msg = NewError("exit", msg).Error()
	_, _ = os.Stderr.Write([]byte(msg))
	os.Exit(code)
}

// WrapError wraps given error into a new error with format.
func WrapError(code, msg string, err error) *Error {
	var wrap string
	if err != nil {
		wrap = delim + err.Error()
	}

	return &Error{Code: code, Msg: msg + wrap}
}
