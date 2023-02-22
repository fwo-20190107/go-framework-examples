package infra

import (
	"examples/model"
	"net/http"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

func (fn HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewLogContext(r.Context())
	r = r.WithContext(ctx)
	defer model.Logger.Send(ctx)

	fn(w, r)
}

var _ http.Handler = (*HttpHandler)(nil)
