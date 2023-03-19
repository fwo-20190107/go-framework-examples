package web

import (
	"examples/errors"
	"examples/internal/http/interface/infra"
	"examples/internal/http/registry"
	"fmt"
	"net/http"
)

type HttpHandler infra.HttpHandler

func (fn HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	httpCtx := newHttpContext(w, r)

	if err := fn(ctx, httpCtx); err != nil {
		status := errors.HTTPStatus(err.Err)
		if err.Err != nil {
			switch {
			case status >= 500:
				registry.Logger.Err(ctx, fmt.Sprint(err.Err))
			case status >= 400:
				registry.Logger.Warn(ctx, fmt.Sprint(err.Err))
			}
		}
		httpCtx.WriteError(status, err.Response)
	}
}

var _ http.Handler = (HttpHandler)(nil)
