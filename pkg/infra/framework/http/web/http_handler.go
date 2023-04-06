package web

import (
	"examples/pkg/adapter/framework/http/infra"
	"examples/pkg/errors"
	"examples/pkg/logger"
	"fmt"
	"net/http"
)

type HttpHandler infra.HttpHandler

func (fn HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	httpCtx := newHttpContext(w, r)

	if err := fn(ctx, httpCtx); err != nil {
		status := errors.HTTPStatus(err.Error)
		if err.Error != nil {
			switch {
			case status >= 500:
				logger.L.Err(fmt.Sprint(err.Error))
			case status >= 400:
				logger.L.Warn(fmt.Sprint(err.Error))
			}
		}
		httpCtx.WriteError(status, err.HTTPError)
	}
}

var _ http.Handler = (HttpHandler)(nil)
