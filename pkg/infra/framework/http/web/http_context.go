package web

import (
	"encoding/json"
	"examples/pkg/adapter/framework/http/infra"
	commonInfra "examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"
	"net/http"
	"net/url"
	"strings"
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

func (c *httpContext) Header() http.Header {
	return c.r.Header
}

func (c *httpContext) Vars(prefix string, keys ...string) (map[string]string, error) {
	// パス指定漏れの対策
	// /paths でも /paths/ でも同じ結果を得られるようにしておく
	path := strings.TrimPrefix(strings.TrimPrefix(c.URL().Path, prefix), "/")
	param := strings.Split(path, "/")
	if len(param) > len(keys) {
		return nil, errors.Errorf(code.CodeBadRequest, "invalid request path: %s", c.URL().Path)
	}

	vars := make(map[string]string, len(param))
	for i, p := range param {
		if len(p) == 0 {
			continue
		}
		vars[keys[i]] = p
	}
	return vars, nil
}

func (c *httpContext) Decode(v any) error {
	decoder := json.NewDecoder(c.r.Body)
	if err := decoder.Decode(&v); err != nil {
		return errors.Wrap(code.CodeBadRequest, err)
	}
	return nil
}

func (c *httpContext) WriteJSON(status int, body any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return errors.Errorf(code.CodeInternal, err.Error())
	}

	c.w.WriteHeader(status)
	if _, err := c.w.Write(jsonBody); err != nil {
		return err
	}
	return nil
}

func (c *httpContext) WriteError(status int, res *commonInfra.HTTPError) error {
	return c.WriteJSON(status, res)
}

var _ infra.HttpContext = (*httpContext)(nil)
