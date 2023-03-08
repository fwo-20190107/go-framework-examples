package handler

import (
	"database/sql"
	"examples/internal/http/interface/infra"
	"examples/internal/http/logic"
	"examples/message"
	"fmt"
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

func (h *userHandler) getUserByID(ctx infra.HttpContext) *infra.HttpError {
	user, err := h.userLogic.GetUserByID(ctx.Context(), 1)
	if err != nil {
		var msg string
		switch err {
		case sql.ErrNoRows:
			msg = fmt.Sprintf("user not found. userID=%d", 1)
		default:
			msg = err.Error()
		}
		return &infra.HttpError{Msg: msg, Code: http.StatusInternalServerError}
	}

	if err := ctx.WriteJSON(http.StatusOK, user); err != nil {
		return &infra.HttpError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}
	return nil
}

func (h *userHandler) HandleRoot(ctx infra.HttpContext) *infra.HttpError {
	switch ctx.Method() {
	case http.MethodGet:
		return h.getUserByID(ctx)
	}
	return message.ErrNotFound
}
