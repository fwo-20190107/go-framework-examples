package web

import (
	"examples/pkg/adapter/framework/chi/infra"
	"examples/pkg/errors"
	"examples/pkg/logger"
	"fmt"
	"net/http"
)

type ChiHandler infra.ChiHandler

func (fn ChiHandler) Exec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	chiCtx := newChiContext(w, r)

	if err := fn(ctx, chiCtx); err != nil {
		status := errors.HTTPStatus(err.Error)
		if err.Error != nil {
			switch {
			case status >= http.StatusInternalServerError:
				logger.L.Err(fmt.Sprint(err.Error))
			case status >= http.StatusBadRequest:
				logger.L.Warn(fmt.Sprint(err.Error))
			}
		}
		chiCtx.WriteError(status, err.HTTPError)
	}
}

var _ http.HandlerFunc = (ChiHandler)(nil).Exec
