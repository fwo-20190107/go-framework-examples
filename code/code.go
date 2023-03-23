package code

import "errors"

type ErrorCode struct {
	err error
}

var (
	ErrOK = newErrorCode("ok")

	ErrBadRequest   = newErrorCode("bad request")
	ErrUnauthorized = newErrorCode("unauthorized")
	ErrValidParam   = newErrorCode("valid parameter")
	ErrNotFound     = newErrorCode("not found")
	ErrOutOfTerm    = newErrorCode("out of term")

	ErrDatabase = newErrorCode("database error")
	ErrInternal = newErrorCode("internal error")

	ErrUnknown = newErrorCode("unknown error")
)

func newErrorCode(msg string) *ErrorCode {
	return &ErrorCode{err: errors.New(msg)}
}

func (e *ErrorCode) Is(err error) bool {
	return errors.Is(e.err, err)
}

func (e *ErrorCode) Error() string {
	return e.err.Error()
}

var _ error = (*ErrorCode)(nil)
