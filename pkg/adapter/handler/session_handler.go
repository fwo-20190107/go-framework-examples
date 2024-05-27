package handler

import (
	"examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/logic"
	"examples/pkg/logic/iodata"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionHandler interface {
	Signin(c *gin.Context) *infra.HandleError
	Signout(c *gin.Context) *infra.HandleError
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

func (h *sessionHandler) Signin(c *gin.Context) *infra.HandleError {
	var input *iodata.SigninInput
	if err := c.ShouldBindJSON(&input); err != nil {
		r := NewHTTPError("クライアントエラー", "リクエストされた値の取得に失敗しました")
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	if err := input.Validate(); err != nil {
		return &infra.HandleError{HTTPError: ErrValidParam, Error: err}
	}

	userID, err := h.sessionLogic.Signin(c, input)
	if err != nil {
		// サインイン失敗時は、後の攻撃を抑制するため詳細のエラーは返却しない
		// e.g. ログインIDが存在しない / パスワードが不一致
		return &infra.HandleError{HTTPError: ErrFailedSignin, Error: err}
	}

	user, err := h.userLogic.GetByID(c, userID)
	if err != nil {
		r := ErrUnexpected
		switch {
		case errors.Is(err, code.CodeNotFound):
			// 正しくユーザー登録が行われていればエラーとならない
			// このケースは問題があるのでエラーレベルを引き上げる
			err = errors.Errorf(code.CodeInternal, err.Error())
			r = NewHTTPError("整合性エラー", "ログインIDに紐付くユーザー情報が見つかりません")
		}
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	token, err := h.sessionLogic.Start(c, user.UserID)
	if err != nil {
		r := NewHTTPError("サーバーエラー", "ログイントークンの発行に失敗しました")
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

func (h *sessionHandler) Signout(c *gin.Context) *infra.HandleError {
	h.sessionLogic.Signout(c)
	return nil
}
