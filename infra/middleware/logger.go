package middleware

import (
	"examples/infra"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func WithLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := infra.NewResponseWriter(w, r)
		next.ServeHTTP(writer, r)

		logger.Info().Object("accesslog", writer).Send()
	}
}

func init() {
	logger = zerolog.New(os.Stdout)
}
