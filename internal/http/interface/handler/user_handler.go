package handler

import (
	"context"
	"examples/code"
	"examples/errors"
	"examples/internal/http/interface/infra"
	"examples/internal/http/logic"
	"net/http"
)

type userHandler struct {
	userLogic logic.UserLogic
}

func NewUserHandler(userLogic logic.UserLogic) *userHandler {
	return &userHandler{
		userLogic: userLogic,
	}
}

func (h *userHandler) getUserByID(ctx context.Context, httpCtx infra.HttpContext) *infra.HttpError {
	user, err := h.userLogic.GetByID(ctx, 1)
	if err != nil {
		r := ErrUnexpected
		switch {
		case errors.Is(err, code.ErrNotFound):
			r = newErrorResponse("エラー", "ユーザーデータなし")
		}
		return &infra.HttpError{Response: r, Err: err}
	}

	if err := httpCtx.WriteJSON(http.StatusOK, user); err != nil {
		return &infra.HttpError{Response: ErrUnexpected, Err: err}
	}
	return nil
}

func (h *userHandler) HandleRoot(ctx context.Context, httpCtx infra.HttpContext) *infra.HttpError {
	switch httpCtx.Method() {
	case http.MethodGet:
		return h.getUserByID(ctx, httpCtx)
	}
	return nil
}
