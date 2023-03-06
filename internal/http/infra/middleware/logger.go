package middleware

import (
	"examples/internal/http/infra/log"
	"examples/internal/http/infra/web"
	"examples/internal/http/registry"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func WithLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := web.NewResponseWriter(w, r)

		ctx := log.NewLogContext(r.Context())
		r = r.WithContext(ctx)

		next.ServeHTTP(writer, r)

		if err := registry.Logger.Send(ctx); err != nil {
			fmt.Println(err)
		}
		logger.Info().Object("accesslog", writer).Send()
	}
}

func init() {
	logger = zerolog.New(os.Stdout)
}
