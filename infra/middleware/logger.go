package middleware

import (
	"examples/infra"
	"examples/model"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func WithLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := infra.NewResponseWriter(w)

		ctx := infra.NewLogContext(r.Context())
		r = r.WithContext(ctx)

		next.ServeHTTP(writer, r)

		model.Logger.Send(ctx)
		logger.Info().Object("accesslog", writer).Send()
	}
}

func init() {
	logger = zerolog.New(os.Stdout)
}
