package middleware

import (
	"examples/internal/http/registry"
	"examples/internal/http/util"
	"net/http"
)

const HEADER_AUTHORIZATION = "Authorization"

func CheckToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(HEADER_AUTHORIZATION)
		if len(token) == 0 {
			next = http.HandlerFunc(unauthorized)
		} else {
			ctx := r.Context()
			userID, ok := registry.SessionManager.Load(ctx, token)
			if !ok {
				next = http.HandlerFunc(unauthorized)
			} else {
				ctx = util.SetUserInfo(ctx, token, userID)
			}
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	}
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}
