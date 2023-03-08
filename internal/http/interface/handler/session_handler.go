package handler

import (
	"examples/internal/http/interface/infra"
	"examples/internal/http/logic"
	"examples/internal/http/logic/iodata"
	"examples/message"
	"net/http"
	"strings"
)

type sessionHandler struct {
	userLogic logic.UserLogic
}

func NewSessionHandler(userLogic logic.UserLogic) *sessionHandler {
	return &sessionHandler{
		userLogic: userLogic,
	}
}

func (h *sessionHandler) signin(ctx infra.HttpContext) *infra.HttpError {
	var in iodata.SigninInput
	if err := ctx.Decode(&in); err != nil {
		return &infra.HttpError{Msg: message.ErrParseForm.Error(), Code: http.StatusBadRequest}
	}

	if err := in.Validate(); err != nil {
		return &infra.HttpError{Msg: err.Error(), Code: http.StatusBadRequest, Err: err}
	}

	if _, err := h.userLogic.Signin(ctx.Context(), in.LoginID, in.Password); err != nil {
		return &infra.HttpError{Msg: "login failed", Code: http.StatusUnauthorized, Err: err}
	}
	return nil
}

func (h *sessionHandler) signout(ctx infra.HttpContext) *infra.HttpError {
	h.userLogic.Signout(ctx.Context())
	return nil
}

func (h *sessionHandler) HandleRoot(ctx infra.HttpContext) *infra.HttpError {
	path := strings.TrimPrefix(ctx.URL().Path, "session/")
	if len(path) > 0 {
		return message.ErrNotFound
	}

	switch ctx.Method() {
	case http.MethodPost:
		return h.signin(ctx)
	case http.MethodDelete:
		return h.signout(ctx)
	}
	return message.ErrNotFound
}
