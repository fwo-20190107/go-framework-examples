package middleware

import (
	"examples/internal/http/registry"
	"fmt"
	"net/http"
	"runtime/debug"
)

func WithRecover(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer func() {
			if err := recover(); err != nil {
				registry.Logger.Fatal(ctx, fmt.Sprintf("Msg: %v\n%s\n", err, string(debug.Stack())))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("fatal error"))
			}
		}()
		next.ServeHTTP(w, r)
	}
}
