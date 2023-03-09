package infra

import (
	"context"
	"net/url"
)

type HttpHandler func(httpCtx HttpContext) *HttpError

type HttpError struct {
	Code int
	Msg  string
	Err  error
}

type HttpContext interface {
	Context() context.Context
	URL() *url.URL
	Method() string
	Decode(v any) error
	WriteJSON(code int, body any) error
	WriteError(code int, msg string) error
}
