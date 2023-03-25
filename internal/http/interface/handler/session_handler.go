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

func (h *sessionHandler) Signin(ctx context.Context, httpCtx infra.HttpContext) *infra.HandleError {
	if httpCtx.Method() != http.MethodPost {
		return &infra.HandleError{HTTPError: ErrPathNotExist}
	}

	var input *iodata.SigninInput
	if err := httpCtx.Decode(&input); err != nil {
		r := NewHTTPError("クライアントエラー", "リクエストされた値の取得に失敗しました")
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	if err := input.Validate(); err != nil {
		return &infra.HandleError{HTTPError: ErrValidParam, Error: err}
	}

	userID, err := h.sessionLogic.Signin(ctx, input)
	if err != nil {
		// サインイン失敗時は、後の攻撃を抑制するため詳細のエラーは返却しない
		// e.g. ログインIDが存在しない / パスワードが不一致
		return &infra.HandleError{HTTPError: ErrFailedSignin, Error: err}
	}

	user, err := h.userLogic.GetByID(ctx, userID)
	if err != nil {
		r := ErrUnexpected
		switch {
		case errors.Is(err, code.ErrNotFound):
			// 正しくユーザー登録が行われていればエラーとならない
			// このケースは問題があるのでエラーレベルを引き上げる
			err = errors.Errorf(code.ErrInternal, err.Error())
			r = NewHTTPError("整合性エラー", "ログインIDに紐付くユーザー情報が見つかりません")
		}
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	token, err := h.sessionLogic.Start(ctx, user.UserID)
	if err != nil {
		r := NewHTTPError("サーバーエラー", "ログイントークンの発行に失敗しました")
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	httpCtx.WriteJSON(http.StatusOK, &iodata.SigninOutput{
		Token:     token,
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	})
	return nil
}

func (h *sessionHandler) Signout(ctx context.Context, httpCtx infra.HttpContext) *infra.HandleError {
	if httpCtx.Method() != http.MethodDelete {
		return &infra.HandleError{HTTPError: ErrPathNotExist}
	}

	h.sessionLogic.Signout(ctx)
	return nil
}
