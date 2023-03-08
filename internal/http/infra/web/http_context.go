package web

import (
	"context"
	"encoding/json"
	"examples/internal/http/interface/infra"
	"net/http"
	"net/url"
)

type httpContext struct {
	w http.ResponseWriter
	r *http.Request
}

func newHttpContext(w http.ResponseWriter, r *http.Request) *httpContext {
	return &httpContext{w: w, r: r}
}

func (c *httpContext) Context() context.Context {
	return c.r.Context()
}

func (c *httpContext) URL() *url.URL {
	return c.r.URL
}

func (c *httpContext) Method() string {
	return c.r.Method
}

func (c *httpContext) Decode(v any) error {
	decoder := json.NewDecoder(c.r.Body)
	if err := decoder.Decode(&v); err != nil {
		return err
	}
	return nil
}

var _ infra.HttpContext = (*httpContext)(nil)
