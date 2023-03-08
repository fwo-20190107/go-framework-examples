package handler

import (
	"context"
	"examples/internal/http/logic"
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

func (h *sessionHandler) signin(ctx context.Context, w http.ResponseWriter, r *http.Request) *handleError {
	if err := r.ParseForm(); err != nil {
		return &handleError{msg: message.ErrParseForm.Error(), code: http.StatusBadRequest}
	}

	if _, err := h.userLogic.Signin(ctx, "admin", "admin"); err != nil {
		return &handleError{msg: "login failed", code: http.StatusUnauthorized}
	}
	return nil
}

func (h *sessionHandler) signout(ctx context.Context, w http.ResponseWriter, r *http.Request) *handleError {
	h.userLogic.Signout(ctx)
	return nil
}

func (h *sessionHandler) HandleRoot(w http.ResponseWriter, r *http.Request) *handleError {
	path := strings.TrimPrefix(r.URL.Path, "session/")
	if len(path) > 0 {
		return ErrNotFound
	}

	ctx := r.Context()
	switch r.Method {
	case http.MethodPost:
		return h.signin(ctx, w, r)
	case http.MethodDelete:
		return h.signout(ctx, w, r)
	default:
		http.NotFound(w, r)
	}
	return nil
}
