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
	userLogic    logic.UserLogic
	sessionLogic logic.SessionLogic
}

func NewSessionHandler(userLogic logic.UserLogic, sessionLogic logic.SessionLogic) *sessionHandler {
	return &sessionHandler{
		userLogic:    userLogic,
		sessionLogic: sessionLogic,
	}
}

func (h *sessionHandler) Signup(ctx context.Context, httpCtx infra.HttpContext) *infra.HttpError {
	return nil
}

func (h *sessionHandler) Signin(ctx context.Context, httpCtx infra.HttpContext) *infra.HttpError {
	if httpCtx.Method() != http.MethodPost {
		return &infra.HttpError{Response: ErrPathNotExist}
	}

	var in iodata.SigninInput
	if err := httpCtx.Decode(&in); err != nil {
		r := newErrorResponse("クライアントエラー", "リクエストされた値の取得に失敗しました")
		return &infra.HttpError{Response: r, Err: err}
	}

	if err := in.Validate(); err != nil {
		return &infra.HttpError{Response: ErrValidParam, Err: err}
	}

	userID, err := h.sessionLogic.Signin(ctx, in.LoginID, in.Password)
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

	token, err := h.sessionLogic.Start(ctx, user.UserID)
	if err != nil {
		r := newErrorResponse("サーバーエラー", "ログイントークンの発行に失敗しました")
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

func (h *sessionHandler) Signout(ctx context.Context, httpCtx infra.HttpContext) *infra.HttpError {
	if httpCtx.Method() != http.MethodDelete {
		return &infra.HttpError{Response: ErrPathNotExist}
	}

	h.sessionLogic.Signout(ctx)
	return nil
}
