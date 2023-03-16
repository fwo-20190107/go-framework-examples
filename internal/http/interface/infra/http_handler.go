package infra

import (
	"context"
	"net/url"
)

type HttpHandler func(ctx context.Context, httpCtx HttpContext) *HttpError

type HttpError struct {
	Code int
	Msg  string
	Err  error
}

type HttpContext interface {
	URL() *url.URL
	Method() string
	Decode(v any) error
	WriteJSON(code int, body any) error
	WriteError(code int, msg string) error
}
