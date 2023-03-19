package infra

import (
	"context"
	"net/url"
)

type HttpHandler func(ctx context.Context, httpCtx HttpContext) *HttpError

type ErrorResponse struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type HttpError struct {
	Response *ErrorResponse
	Err      error
}

type HttpContext interface {
	URL() *url.URL
	Method() string
	Decode(v any) error
	WriteJSON(code int, body any) error
	WriteError(code int, res *ErrorResponse) error
}
