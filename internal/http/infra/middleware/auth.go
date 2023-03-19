package middleware

import (
	"context"
	"examples/code"
	"examples/errors"
	"examples/internal/http/interface/infra"
	"examples/internal/http/registry"
	"examples/internal/http/util"
	"net/http"
)

const HEADER_AUTHORIZATION = "Authorization"

func CheckToken(next infra.HttpHandler) infra.HttpHandler {
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
			userID, ok := registry.SessionManager.Load(ctx, token)
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

func unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}
