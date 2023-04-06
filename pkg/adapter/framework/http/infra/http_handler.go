package infra

import (
	"context"
	"examples/pkg/adapter/infra"
	"net/http"
	"net/url"
)

type HttpHandler func(ctx context.Context, httpCtx HttpContext) *infra.HandleError

type HttpContext interface {
	URL() *url.URL
	Method() string
	Header() http.Header
	Vars(prefix string, keys ...string) (map[string]string, error)
	Decode(v any) error
	WriteJSON(code int, body any) error
	WriteError(code int, res *infra.HTTPError) error
}
