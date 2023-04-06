package handler

import (
	"context"
	"examples/pkg/adapter/handler"
	"examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/logic"
	"examples/pkg/logic/iodata"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SessionHandler interface {
	Signin(ctx context.Context, c echo.Context) *infra.HandleError
	Signout(ctx context.Context, c echo.Context) *infra.HandleError
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

func (h *sessionHandler) Signin(ctx context.Context, c echo.Context) *infra.HandleError {
	var input *iodata.SigninInput
	if err := c.Bind(&input); err != nil {
		r := handler.NewHTTPError("クライアントエラー", "リクエストされた値の取得に失敗しました")
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	if err := input.Validate(); err != nil {
		return &infra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}

	userID, err := h.sessionLogic.Signin(ctx, input)
	if err != nil {
		// サインイン失敗時は、後の攻撃を抑制するため詳細のエラーは返却しない
		// e.g. ログインIDが存在しない / パスワードが不一致
		return &infra.HandleError{HTTPError: handler.ErrFailedSignin, Error: err}
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
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	token, err := h.sessionLogic.Start(ctx, user.UserID)
	if err != nil {
		r := handler.NewHTTPError("サーバーエラー", "ログイントークンの発行に失敗しました")
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	c.JSON(http.StatusOK, &iodata.SigninOutput{
		Token: token,
		UserOutput: iodata.UserOutput{
			UserID:    user.UserID,
			Name:      user.Name,
			Authority: user.Authority,
		},
	})
	return nil
}

func (h *sessionHandler) Signout(ctx context.Context, c echo.Context) *infra.HandleError {
	h.sessionLogic.Signout(ctx)
	return nil
}
