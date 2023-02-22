package infra

import (
	"net/http"

	"github.com/rs/zerolog"
)

type responseWriter struct {
	w http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w}
}

func (w *responseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *responseWriter) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

func (w *responseWriter) WriteHeader(status int) {
	w.w.WriteHeader(status)
}

func (w *responseWriter) MarshalZerologObject(e *zerolog.Event) {
}

var _ http.ResponseWriter = (*responseWriter)(nil)
var _ zerolog.LogObjectMarshaler = (*responseWriter)(nil)
