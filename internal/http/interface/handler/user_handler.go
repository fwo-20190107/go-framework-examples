package handler

import (
	"context"
	"database/sql"
	"examples/internal/http/logic"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

type userHandler struct {
	userLogic logic.UserLogic
}

func NewUserHandler(userLogic logic.UserLogic) *userHandler {
	return &userHandler{
		userLogic: userLogic,
	}
}

func (h *userHandler) getUserByID(ctx context.Context, w http.ResponseWriter, r *http.Request) *handleError {
	user, err := h.userLogic.GetUserByID(ctx, 1)
	if err != nil {
		var msg string
		switch err {
		case sql.ErrNoRows:
			msg = fmt.Sprintf("user not found. userID=%d", 1)
		default:
			msg = err.Error()
		}
		return &handleError{msg: msg, code: http.StatusInternalServerError}
	}

	bytesUser, err := json.Marshal(user)
	if err != nil {
		return &handleError{msg: err.Error(), code: http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytesUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (h *userHandler) HandleRoot(w http.ResponseWriter, r *http.Request) (err *handleError) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		err = h.getUserByID(ctx, w, r)
	default:
		http.NotFound(w, r)
	}
	return
}
