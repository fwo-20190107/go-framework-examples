package errors

import (
	"examples/pkg/code"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type applicationError struct {
	code code.ErrorCode
	err  error
}

func Errorf(c code.ErrorCode, format string, args ...any) error {
	if errors.Is(c, code.CodeOK) {
		return nil
	}
	return &applicationError{
		code: c,
		err:  fmt.Errorf(format, args...),
		// err:  errors.Errorf(format, args...),
	}
}

func Wrap(c code.ErrorCode, err error) error {
	if err == nil || errors.Is(c, code.CodeOK) {
		return nil
	}
	return &applicationError{
		code: c,
		err:  err,
		// err:  errors.WithStack(err),
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

func errorCode(err error) code.ErrorCode {
	if err == nil {
		return code.CodeOK
	}
	var e *applicationError
	if errors.As(err, &e) {
		return e.code
	}
	return code.CodeUnknown
}

var _ error = (*applicationError)(nil)

var statusMap = map[code.ErrorCode]int{
	code.CodeOK: http.StatusOK,

	code.CodeBadRequest:   http.StatusBadRequest,
	code.CodeValidParam:   http.StatusBadRequest,
	code.CodeUnauthorized: http.StatusUnauthorized,
	code.CodeNotFound:     http.StatusNotFound,
	code.CodeOutOfTerm:    http.StatusNotFound,

	code.CodeDatabase: http.StatusInternalServerError,
	code.CodeInternal: http.StatusInternalServerError,

	code.CodeUnknown: http.StatusInternalServerError,
}
