package middleware

import (
	"context"
	"examples/code"
	"examples/errors"
	"examples/pkg/http/logic/repository"
	"examples/pkg/http/registry"
	"examples/pkg/http/util"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const tokenPrefix = "Bearer "

type AuthMiddleware struct {
	sessionRepository repository.SessionRepository
}

func NewAuthMiddleware(sessionRepository repository.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{sessionRepository: sessionRepository}
}

func (m *AuthMiddleware) WithCheckToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, tokenPrefix) {
			warnLog(ctx, errors.Errorf(code.ErrUnauthorized, "Bearer token is required."))
			unauthorized(w)
			return
		}
		token = strings.TrimPrefix(token, tokenPrefix)

		userID, ok := m.sessionRepository.Get(ctx, token)
		if !ok {
			warnLog(ctx, errors.Errorf(code.ErrUnauthorized, "Request token is invalid."))
			unauthorized(w)
			return
		}
		ctx = util.SetUserInfo(ctx, token, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func warnLog(ctx context.Context, err error) {
	if err := registry.Logger.Warn(ctx, fmt.Sprint(err)); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("invalid token"))
}
