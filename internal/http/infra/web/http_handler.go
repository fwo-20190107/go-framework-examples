package web

import (
	"examples/internal/http/interface/infra"
	"net/http"
)

type HttpHandler infra.HttpHandler

func (fn HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := newHttpContext(w, r)

	if err := fn(ctx); err != nil {
		http.Error(w, err.Msg, err.Code)
	}
}

var _ http.Handler = (HttpHandler)(nil)
