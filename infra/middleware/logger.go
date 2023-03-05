package middleware

import (
	"examples/entity"
	"examples/infra"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func WithLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := infra.NewResponseWriter(w, r)

		ctx := infra.NewLogContext(r.Context())
		r = r.WithContext(ctx)

		next.ServeHTTP(writer, r)

		if err := entity.Logger.Send(ctx); err != nil {
			fmt.Println(err)
		}
		logger.Info().Object("accesslog", writer).Send()
	}
}

func init() {
	logger = zerolog.New(os.Stdout)
}
