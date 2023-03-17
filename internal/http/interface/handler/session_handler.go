package handler

import (
	"context"
	"examples/internal/http/interface/infra"
	"examples/internal/http/logic"
	"examples/internal/http/logic/iodata"
	"examples/message"
	"net/http"
	"strings"
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
		return &infra.HttpError{Msg: message.ErrParseForm.Error(), Code: http.StatusBadRequest}
	}

	if err := in.Validate(); err != nil {
		return &infra.HttpError{Msg: err.Error(), Code: http.StatusBadRequest, Err: err}
	}

	userID, err := h.loginLogic.Signin(ctx, in.LoginID, in.Password)
	if err != nil {
		return &infra.HttpError{Msg: "login failed", Code: http.StatusUnauthorized, Err: err}
	}

	user, err := h.userLogic.GetByID(ctx, userID)
	if err != nil {
		return &infra.HttpError{Msg: "login user is not found", Code: http.StatusInternalServerError, Err: err}
	}

	token, err := h.loginLogic.PublishToken(ctx, user.UserID)
	if err != nil {
		return &infra.HttpError{Msg: "failed to publish token", Code: http.StatusInternalServerError, Err: err}
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
	path := strings.TrimPrefix(httpCtx.URL().Path, "/session")
	if len(path) > 0 {
		return message.ErrNotFound
	}

	switch httpCtx.Method() {
	case http.MethodPost:
		return h.signin(ctx, httpCtx)
	case http.MethodDelete:
		return h.signout(ctx, httpCtx)
	}
	return message.ErrNotFound
}
