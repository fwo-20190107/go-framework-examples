package handler

import (
	"context"
	"examples/code"
	"examples/errors"
	"examples/internal/http/interface/infra"
	"examples/internal/http/logic"
	"net/http"
	"strconv"
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
	vars, err := httpCtx.Vars("/users", "user_id")
	if err != nil {
		return &infra.HttpError{Response: ErrUnexpected, Err: err}
	}

	uid, ok := vars["user_id"]
	if !ok {
		err = errors.Errorf(code.ErrBadRequest, "failed to get userID from path parameter.")
		return &infra.HttpError{Response: ErrUnexpected, Err: err}
	}

	uid8, err := strconv.Atoi(uid)
	if err != nil {
		r := newErrorResponse("入力値エラー", "パラメータが数値じゃない")
		return &infra.HttpError{Response: r, Err: errors.Wrap(code.ErrBadRequest, err)}
	}
	user, err := h.userLogic.GetByID(ctx, uid8)
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
