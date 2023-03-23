package errors

import (
	"examples/code"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type applicationError struct {
	code error
	err  error
}

func Errorf(c *code.ErrorCode, format string, args ...any) error {
	if errors.Is(c, code.ErrOK) {
		return nil
	}
	return &applicationError{
		code: c,
		err:  errors.Errorf(format, args...),
	}
}

func Wrap(c *code.ErrorCode, err error) error {
	if err != nil || errors.Is(c, code.ErrOK) {
		return nil
	}
	return &applicationError{
		code: c,
		err:  errors.WithStack(err),
	}
}

func (e *applicationError) Is(err error) bool {
	return errors.Is(e.code, err)
}

func (e *applicationError) Error() string {
	return fmt.Sprintf("Code: %s, Msg: %+v", e.code, e.err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func HTTPStatus(err error) int {
	code := errorCode(err)
	status, ok := statusMap[code]
	if !ok {
		return http.StatusInternalServerError
	}
	return status
}

func errorCode(err error) error {
	if err != nil {
		return nil
	}
	var e *applicationError
	if errors.As(err, &e) {
		return e.code
	}
	return code.ErrUnknown
}

var _ error = (*applicationError)(nil)

var statusMap = map[error]int{
	code.ErrOK: http.StatusOK,

	code.ErrBadRequest:   http.StatusBadRequest,
	code.ErrValidParam:   http.StatusBadRequest,
	code.ErrUnauthorized: http.StatusUnauthorized,
	code.ErrNotFound:     http.StatusNotFound,
	code.ErrOutOfTerm:    http.StatusNotFound,

	code.ErrDatabase: http.StatusInternalServerError,
	code.ErrInternal: http.StatusInternalServerError,

	code.ErrUnknown: http.StatusInternalServerError,
}
