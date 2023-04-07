package middleware

import (
	"examples/pkg/logger"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type formatter struct {
	logger zerolog.Logger
}

type entry struct {
	logger zerolog.Logger
}

func WithLogger(w io.Writer) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&formatter{logger: zerolog.New(w)})
}

func (f *formatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	logger := f.logger.With().
		Str("requestURI", r.RequestURI).
		Str("method", r.Method).
		Int64("size", r.ContentLength).
		Str("remoteIP", r.RemoteAddr).
		Str("agent", r.UserAgent()).
		Logger()

	return &entry{logger: logger}
}

func (e *entry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	e.logger.Info().
		Int("status", status).
		Str("latency", elapsed.String()).
		Send()
}

func (m *entry) Panic(v interface{}, stack []byte) {
	logger.L.Fatal(fmt.Sprintf("%s: %s", fmt.Sprint(v), string(stack)))
}

var _ middleware.LogFormatter = (*formatter)(nil)
var _ middleware.LogEntry = (*entry)(nil)
