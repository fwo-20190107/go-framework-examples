package middleware

import (
	"context"
	"examples/code"
	"examples/errors"
	"examples/internal/http/interface/infra"
	"examples/internal/http/logic/repository"
	"examples/internal/http/util"
)

const HEADER_AUTHORIZATION = "Authorization"

type AuthMiddleware struct {
	sessionRepository repository.SessionRepository
}

func NewAuthMiddleware(sessionRepository repository.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{sessionRepository: sessionRepository}
}

func (m *AuthMiddleware) CheckToken(next infra.HttpHandler) infra.HttpHandler {
	return func(ctx context.Context, httpCtx infra.HttpContext) *infra.HttpError {
		token := httpCtx.Header().Get(HEADER_AUTHORIZATION)
		if len(token) == 0 {
			err := errors.Errorf(code.ErrUnauthorized, "access token is not set")
			r := &infra.ErrorResponse{
				Title: "認証エラー",
				Body:  "トークンは必須です",
			}
			return &infra.HttpError{Response: r, Err: err}
		} else {
			userID, ok := m.sessionRepository.Get(ctx, token)
			if !ok {
				err := errors.Errorf(code.ErrUnauthorized, "illegal token")
				r := &infra.ErrorResponse{
					Title: "認証エラー",
					Body:  "不正なトークンです",
				}
				return &infra.HttpError{Response: r, Err: err}
			}
			ctx = util.SetUserInfo(ctx, token, userID)
		}
		return next(ctx, httpCtx)
	}
}
