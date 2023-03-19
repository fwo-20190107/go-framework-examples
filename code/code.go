package code

import "errors"

type ErrorCode struct {
	err error
}

var (
	ErrOK = new("ok")

	ErrBadRequest   = new("bad request")
	ErrUnauthorized = new("unauthorized")
	ErrValidParam   = new("valid parameter")
	ErrNotFound     = new("not found")
	ErrOutOfTerm    = new("out of term")

	ErrDatabase = new("database error")
	ErrInternal = new("internal error")

	ErrUnknown = new("unknown error")
)

func new(msg string) *ErrorCode {
	return &ErrorCode{err: errors.New(msg)}
}

func (e *ErrorCode) Is(err error) bool {
	return errors.Is(e.err, err)
}

func (e *ErrorCode) Error() string {
	return e.err.Error()
}

var _ error = (*ErrorCode)(nil)
