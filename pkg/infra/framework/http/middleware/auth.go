package middleware

import (
	"context"
	"examples/pkg/adapter/infra"
	"examples/pkg/adapter/repository"
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/logger"
	IRepository "examples/pkg/logic/repository"
	"examples/pkg/util"
	"fmt"
	"net/http"
	"strings"
)

const tokenPrefix = "Bearer "

var AuthMw *AuthMiddleware

type AuthMiddleware struct {
	sessionRepository IRepository.SessionRepository
}

func InitAuthMiddleware(store infra.LocalStore) {
	AuthMw = &AuthMiddleware{sessionRepository: repository.NewSessionRepository(store)}
}

func (m *AuthMiddleware) WithCheckToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, tokenPrefix) {
			warnLog(ctx, errors.Errorf(code.CodeUnauthorized, "Bearer token is required."))
			unauthorized(w)
			return
		}
		token = strings.TrimPrefix(token, tokenPrefix)

		userID, ok := m.sessionRepository.Get(ctx, token)
		if !ok {
			warnLog(ctx, errors.Errorf(code.CodeUnauthorized, "Request token is invalid."))
			unauthorized(w)
			return
		}
		ctx = util.SetUserInfo(ctx, token, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func warnLog(ctx context.Context, err error) {
	logger.L.Warn(fmt.Sprint(err))
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("invalid token"))
}
