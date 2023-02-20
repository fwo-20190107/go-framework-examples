package infra

import (
	"net/http"

	"github.com/rs/zerolog"
)

type responseWriter struct {
	w http.ResponseWriter
	r *http.Request
}

func NewResponseWriter(w http.ResponseWriter, r *http.Request) *responseWriter {
	return &responseWriter{
		w: w,
		r: r,
	}
}

func (r *responseWriter) Header() http.Header {
	return r.w.Header()
}

func (r *responseWriter) Write(p []byte) (int, error) {
	return r.w.Write(p)
}

func (r *responseWriter) WriteHeader(status int) {
	r.w.WriteHeader(status)
}

func (r *responseWriter) MarshalZerologObject(e *zerolog.Event) {
}

var _ http.ResponseWriter = (*responseWriter)(nil)
var _ zerolog.LogObjectMarshaler = (*responseWriter)(nil)
