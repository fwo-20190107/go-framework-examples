package handler

import (
	"database/sql"
	"examples/logic"
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

func (h *userHandler) getUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.userLogic.GetUserByID(ctx, 1)
	if err != nil {
		var msg string
		switch err {
		case sql.ErrNoRows:
			msg = fmt.Sprintf("user not found. userID=%d", 1)
		default:
			msg = err.Error()
		}
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	bytesUser, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytesUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUserByID(w, r)
	default:
		http.NotFound(w, r)
	}
}

var _ http.Handler = (*userHandler)(nil)
