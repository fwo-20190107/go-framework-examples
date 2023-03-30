package middleware

import (
	"examples/pkg/logger"
	"fmt"
	"net/http"
	"runtime/debug"
)

func WithRecover(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer func() {
			if err := recover(); err != nil {
				logger.L.Fatal(ctx, fmt.Sprintf("Msg: %v\n%s\n", err, string(debug.Stack())))
				http.Error(w, "Sorry. Request was interrupted.", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
