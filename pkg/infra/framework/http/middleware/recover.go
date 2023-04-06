package middleware

import (
	"examples/pkg/logger"
	"fmt"
	"net/http"
	"runtime/debug"
)

func WithRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.L.Fatal(fmt.Sprintf("Msg: %v\n%s\n", err, string(debug.Stack())))
				http.Error(w, "Sorry. Request was interrupted.", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
