package middleware

import (
	"context"
	"encoding/binary"
	"examples/model"
	"net/http"

	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"
)

const HEADER_AUTHORIZATION = "Authorization"

var flaker *sonyflake.Sonyflake
var session = make(map[string]int)

func CheckToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(HEADER_AUTHORIZATION)
		if len(token) == 0 {
			next = http.HandlerFunc(unauthorized)
		} else {
			userID, ok := session[token]
			if !ok {
				next = http.HandlerFunc(unauthorized)
			} else {
				ctx := model.SetUserInfo(r.Context(), token, userID)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	}
}

func NewToken() (string, error) {
	id, err := flaker.NextID()
	if err != nil {
		return "", err
	}

	bytesID := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(bytesID, id)

	token, err := bcrypt.GenerateFromPassword(bytesID, 4)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func AddSession(token string, userID int) {
	session[token] = userID
}

func RemoveToken(ctx context.Context) {
	token, err := model.GetAccessToken(ctx)
	if err != nil {
		model.Logger.Warn(ctx, err.Error())
	} else {
		delete(session, token)
	}
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func init() {
	flaker = sonyflake.NewSonyflake(sonyflake.Settings{})
}
