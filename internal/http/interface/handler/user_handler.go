package handler

import (
	"context"
	"examples/code"
	"examples/errors"
	"examples/internal/http/interface/infra"
	"examples/internal/http/logic"
	"examples/internal/http/logic/iodata"
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

func (h *userHandler) Signup(ctx context.Context, httpCtx infra.HttpContext) *infra.HandleError {
	if httpCtx.Method() != http.MethodPost {
		return &infra.HandleError{HTTPError: ErrPathNotExist}
	}

	var input *iodata.SignupInput
	if err := httpCtx.Decode(&input); err != nil {
		return &infra.HandleError{HTTPError: ErrValidParam, Error: err}
	}

	if err := input.Validate(); err != nil {
		return &infra.HandleError{HTTPError: ErrValidParam, Error: err}
	}

	if err := h.userLogic.Signup(ctx, input); err != nil {
		return &infra.HandleError{HTTPError: ErrUnexpected, Error: err}
	}
	return nil
}

func (h *userHandler) getUserByID(ctx context.Context, httpCtx infra.HttpContext) *infra.HandleError {
	vars, err := httpCtx.Vars("/users", "user_id")
	if err != nil {
		return &infra.HandleError{HTTPError: ErrUnexpected, Error: err}
	}

	uid, ok := vars["user_id"]
	if !ok {
		err = errors.Errorf(code.ErrBadRequest, "failed to get userID from path parameter.")
		return &infra.HandleError{HTTPError: ErrUnexpected, Error: err}
	}

	uid8, err := strconv.Atoi(uid)
	if err != nil {
		r := NewHTTPError("入力値エラー", "パラメータが数値じゃない")
		return &infra.HandleError{HTTPError: r, Error: errors.Wrap(code.ErrBadRequest, err)}
	}
	user, err := h.userLogic.GetByID(ctx, uid8)
	if err != nil {
		r := ErrUnexpected
		switch {
		case errors.Is(err, code.ErrNotFound):
			r = NewHTTPError("エラー", "ユーザーデータなし")
		}
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	if err := httpCtx.WriteJSON(http.StatusOK, user); err != nil {
		return &infra.HandleError{HTTPError: ErrUnexpected, Error: err}
	}
	return nil
}

func (h *userHandler) HandleRoot(ctx context.Context, httpCtx infra.HttpContext) *infra.HandleError {
	switch httpCtx.Method() {
	case http.MethodGet:
		return h.getUserByID(ctx, httpCtx)
	}
	return nil
}
