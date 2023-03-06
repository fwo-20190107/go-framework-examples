package handler

import (
	"net/http"
)

type handleError struct {
	code int
	msg  string
}

type AppHandler func(w http.ResponseWriter, r *http.Request) *handleError

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.msg, err.code)
	}
}

var _ http.Handler = (AppHandler)(nil)
