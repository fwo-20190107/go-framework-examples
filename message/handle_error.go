package message

import (
	"examples/internal/http/interface/infra"
	"net/http"
)

var (
	ErrNotFound = &infra.HttpError{
		Code: http.StatusNotFound,
		Msg:  http.StatusText(http.StatusNotFound),
	}
	ErrUnexpected = &infra.HttpError{
		Code: http.StatusInternalServerError,
		Msg:  "error unexpected",
	}
)
