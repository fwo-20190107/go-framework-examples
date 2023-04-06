package handler

import (
	"context"
	"examples/pkg/adapter/framework/http/infra"
	"examples/pkg/adapter/handler"
	cInfra "examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/logic"
	"examples/pkg/logic/iodata"
	"net/http"
)

type SessionHandler interface {
	Signin(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError
	Signout(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError
}

type sessionHandler struct {
	userLogic    logic.UserLogic
	sessionLogic logic.SessionLogic
}

func NewSessionHandler(userLogic logic.UserLogic, sessionLogic logic.SessionLogic) SessionHandler {
	return &sessionHandler{
		userLogic:    userLogic,
		sessionLogic: sessionLogic,
	}
}

// signin godoc
//	@Summary		Sign in of the application.
//	@Description	Password authentication is performed and the issued token is returned.
//	@Tags			session
//	@Accept			json
//	@Produce		json
//	@Param			input	body		iodata.SigninInput	true	"foo"
//	@Success		200		{object}	iodata.SigninOutput
//	@Failure		400		{object}	infra.HTTPError
//	@Failure		401		{object}	infra.HTTPError
//	@Failure		404		{object}	infra.HTTPError
//	@Failure		500		{object}	infra.HTTPError
//	@Router			/signin [post]
func (h *sessionHandler) Signin(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError {
	if httpCtx.Method() != http.MethodPost {
		return &cInfra.HandleError{HTTPError: handler.ErrPathNotExist}
	}

	var input *iodata.SigninInput
	if err := httpCtx.Decode(&input); err != nil {
		r := handler.NewHTTPError("クライアントエラー", "リクエストされた値の取得に失敗しました")
		return &cInfra.HandleError{HTTPError: r, Error: err}
	}

	if err := input.Validate(); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}

	userID, err := h.sessionLogic.Signin(ctx, input)
	if err != nil {
		// サインイン失敗時は、後の攻撃を抑制するため詳細のエラーは返却しない
		// e.g. ログインIDが存在しない / パスワードが不一致
		return &cInfra.HandleError{HTTPError: handler.ErrFailedSignin, Error: err}
	}

	user, err := h.userLogic.GetByID(ctx, userID)
	if err != nil {
		r := handler.ErrUnexpected
		switch {
		case errors.Is(err, code.CodeNotFound):
			// 正しくユーザー登録が行われていればエラーとならない
			// このケースは問題があるのでエラーレベルを引き上げる
			err = errors.Errorf(code.CodeInternal, err.Error())
			r = handler.NewHTTPError("整合性エラー", "ログインIDに紐付くユーザー情報が見つかりません")
		}
		return &cInfra.HandleError{HTTPError: r, Error: err}
	}

	token, err := h.sessionLogic.Start(ctx, user.UserID)
	if err != nil {
		r := handler.NewHTTPError("サーバーエラー", "ログイントークンの発行に失敗しました")
		return &cInfra.HandleError{HTTPError: r, Error: err}
	}

	httpCtx.WriteJSON(http.StatusOK, &iodata.SigninOutput{
		Token: token,
		UserOutput: iodata.UserOutput{
			UserID:    user.UserID,
			Name:      user.Name,
			Authority: user.Authority,
		},
	})
	return nil
}

// signout godoc
//	@Summary		Sign out of the application.
//	@Description	Discard session.
//	@Tags			session
//	@Success		200
//	@Failure		401	{object}	infra.HTTPError
//	@Failure		404	{object}	infra.HTTPError
//	@Failure		500	{object}	infra.HTTPError
//	@Security		Bearer
//	@Router			/signout [delete]
func (h *sessionHandler) Signout(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError {
	if httpCtx.Method() != http.MethodDelete {
		return &cInfra.HandleError{HTTPError: handler.ErrPathNotExist}
	}

	h.sessionLogic.Signout(ctx)
	return nil
}
