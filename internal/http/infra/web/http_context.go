package web

import (
	"encoding/json"
	"examples/code"
	"examples/errors"
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

func (c *httpContext) WriteJSON(status int, body any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return errors.Errorf(code.ErrInternal, err.Error())
	}

	c.w.WriteHeader(status)
	if _, err := c.w.Write(jsonBody); err != nil {
		return err
	}
	return nil
}

func (c *httpContext) WriteError(status int, res *infra.ErrorResponse) error {
	return c.WriteJSON(status, res)
}

var _ infra.HttpContext = (*httpContext)(nil)
