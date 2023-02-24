package handler

import (
	"context"
	"examples/logic"
	"net/http"
)

type sessionHandler struct {
	userLogic logic.UserLogic
}

func NewSessionHandler(userLogic logic.UserLogic) *sessionHandler {
	return &sessionHandler{
		userLogic: userLogic,
	}
}

func (h *sessionHandler) login(ctx context.Context, w http.ResponseWriter, r *http.Request) *handleError {
	if _, err := h.userLogic.Login(ctx, "admin", "admin"); err != nil {
		return &handleError{msg: "login failed", code: http.StatusUnauthorized}
	}
	return nil
}

func (h *sessionHandler) logout(ctx context.Context) {
	h.userLogic.Logout(ctx)
}

func (h *sessionHandler) HandleRoot(w http.ResponseWriter, r *http.Request) (err *handleError) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodPost:
		err = h.login(ctx, w, r)
	case http.MethodDelete:
		h.logout(ctx)
	default:
		http.NotFound(w, r)
	}
	return
}
