package web

import (
	"encoding/json"
	"examples/pkg/adapter/framework/chi/infra"
	cInfra "examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type chiContext struct {
	w http.ResponseWriter
	r *http.Request
}

func newChiContext(w http.ResponseWriter, r *http.Request) *chiContext {
	return &chiContext{w: w, r: r}
}

func (c *chiContext) URLParam(name string) string {
	return chi.URLParam(c.r, name)
}

func (c *chiContext) Decode(v any) error {
	decoder := json.NewDecoder(c.r.Body)
	if err := decoder.Decode(&v); err != nil {
		return errors.Wrap(code.CodeBadRequest, err)
	}
	return nil
}

func (c *chiContext) WriteJSON(status int, body any) error {
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

func (c *chiContext) WriteError(status int, res *cInfra.HTTPError) error {
	return c.WriteJSON(status, res)
}

var _ infra.ChiContext = (*chiContext)(nil)
