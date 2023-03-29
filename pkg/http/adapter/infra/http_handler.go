package infra

import (
	"context"
	"net/http"
	"net/url"
)

type HttpHandler func(ctx context.Context, httpCtx HttpContext) *HandleError

type HTTPError struct {
	Title string `json:"title"`
	Body  any    `json:"body"`
}

type HandleError struct {
	HTTPError *HTTPError
	Error     error
}

type HttpContext interface {
	URL() *url.URL
	Method() string
	Header() http.Header
	Vars(prefix string, keys ...string) (map[string]string, error)
	Decode(v any) error
	WriteJSON(code int, body any) error
	WriteError(code int, res *HTTPError) error
}
