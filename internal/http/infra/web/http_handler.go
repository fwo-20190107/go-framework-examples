package web

import (
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
		if err.Err != nil {
			registry.Logger.Err(ctx, fmt.Sprint(err.Err))
		}
		httpCtx.WriteError(err.Code, err.Msg)
	}
}

var _ http.Handler = (HttpHandler)(nil)
