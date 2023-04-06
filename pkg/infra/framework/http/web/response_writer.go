package web

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type responseWriter struct {
	writer http.ResponseWriter
	req    *http.Request
	start  time.Time
	status int
}

func NewResponseWriter(w http.ResponseWriter, r *http.Request) *responseWriter {
	return &responseWriter{
		writer: w,
		req:    r,
		start:  time.Now(),
	}
}

func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *responseWriter) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.writer.WriteHeader(status)
}

func (w *responseWriter) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("logTime", w.start.Format(time.DateTime)).
		Str("requestUrl", w.req.URL.String()).
		Str("method", w.req.Method).
		Int64("size", w.req.ContentLength).
		Int("status", w.status).
		Str("latency", time.Since(w.start).String()).
		Str("remoteIp", w.req.RemoteAddr).
		Str("agent", w.req.UserAgent())
}

var _ http.ResponseWriter = (*responseWriter)(nil)
var _ zerolog.LogObjectMarshaler = (*responseWriter)(nil)
