package handler

import "net/http"

var (
	ErrNotFound = &handleError{
		code: http.StatusNotFound,
		msg:  http.StatusText(http.StatusNotFound),
	}
)
