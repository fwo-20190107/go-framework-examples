package handler

import (
	"context"
	"examples/code"
	"examples/errors"
	"examples/internal/http/interface/infra"
	"examples/internal/http/logic"
	"examples/internal/http/logic/iodata"
	"net/http"
)

type sessionHandler struct {
	userLogic  logic.UserLogic
	loginLogic logic.LoginLogic
}

func NewSessionHandler(userLogic logic.UserLogic, loginLogic logic.LoginLogic) *sessionHandler {
	return &sessionHandler{
		userLogic:  userLogic,
		loginLogic: loginLogic,
	}
}

func (h *sessionHandler) signin(ctx context.Context, httpCtx infra.HttpContext) *infra.HttpError {
	var in iodata.SigninInput
	if err := httpCtx.Decode(&in); err != nil {
		r := newErrorResponse("クライアントエラー", "リクエストされた値の取得に失敗しました")
		return &infra.HttpError{Response: r, Err: err}
	}

	if err := in.Validate(); err != nil {
		return &infra.HttpError{Response: ErrValidParam, Err: err}
	}

	userID, err := h.loginLogic.Signin(ctx, in.LoginID, in.Password)
	if err != nil {
		// サインイン失敗時は、後の攻撃を抑制するため詳細のエラーは返却しない
		// e.g. ログインIDが存在しない / パスワードが不一致
		return &infra.HttpError{Response: ErrFailedSignin, Err: err}
	}

	user, err := h.userLogic.GetByID(ctx, userID)
	if err != nil {
		r := ErrUnexpected
		switch {
		case errors.Is(err, code.ErrNotFound):
			// 正しくユーザー登録が行われていればエラーとならない
			// このケースは問題があるのでエラーレベルを引き上げる
			err = errors.Errorf(code.ErrInternal, err.Error())
			r = newErrorResponse("整合性エラー", "ログインIDに紐付くユーザー情報が見つかりません")
		}
		return &infra.HttpError{Response: r, Err: err}
	}

	token, err := h.loginLogic.PublishToken(ctx, user.UserID)
	if err != nil {
		r := newErrorResponse("サーバーエラー", "ログイントークンが発行されませんでした")
		return &infra.HttpError{Response: r, Err: err}
	}

	httpCtx.WriteJSON(http.StatusOK, &iodata.SigninOutput{
		Token:     token,
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	})
	return nil
}

func (h *sessionHandler) signout(ctx context.Context, httpCtx infra.HttpContext) *infra.HttpError {
	h.loginLogic.Signout(ctx)
	return nil
}

func (h *sessionHandler) HandleRoot(ctx context.Context, httpCtx infra.HttpContext) *infra.HttpError {
	switch httpCtx.Method() {
	case http.MethodPost:
		return h.signin(ctx, httpCtx)
	case http.MethodDelete:
		return h.signout(ctx, httpCtx)
	}
	return nil
}
