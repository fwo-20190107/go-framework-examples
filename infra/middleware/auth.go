package middleware

import (
	"encoding/binary"
	"net/http"

	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"
)

const HEADER_AUTHORIZATION = "Authorization"

var flaker *sonyflake.Sonyflake
var session map[string]int

func CheckToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(HEADER_AUTHORIZATION)
		if len(token) == 0 {
			next = http.HandlerFunc(unauthorized)
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

func AddToken(userID int, token string) {
	session[token] = userID
}

func RemoveToken(token string) {
	delete(session, token)
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "403 unauthorized", http.StatusUnauthorized)
}

func init() {
	flaker = sonyflake.NewSonyflake(sonyflake.Settings{})
}
