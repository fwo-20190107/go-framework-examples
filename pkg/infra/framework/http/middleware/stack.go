package middleware

import "net/http"

type middleware func(http.Handler) http.Handler

type MiddlewareStack []middleware

func NewMiddlewareStack(mws ...middleware) MiddlewareStack {
	return MiddlewareStack(append([]middleware(nil), mws...))
}

func (s MiddlewareStack) Append(mw middleware) MiddlewareStack {
	return MiddlewareStack(append(s, mw))
}

func (s MiddlewareStack) Then(h http.Handler) http.Handler {
	for i := range s {
		h = s[len(s)-1-i](h)
	}
	return h
}
