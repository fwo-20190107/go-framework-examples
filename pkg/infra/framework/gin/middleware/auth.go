package middleware

import (
	"context"
	"examples/pkg/adapter/handler"
	"examples/pkg/adapter/infra"
	"examples/pkg/adapter/repository"
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/logger"
	IRepository "examples/pkg/logic/repository"
	"examples/pkg/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const tokenPrefix = "Bearer "

var AuthMw *AuthMiddleware

type AuthMiddleware struct {
	sessionRepository IRepository.SessionRepository
}

func InitAuthMiddleware(store infra.LocalStore) {
	AuthMw = &AuthMiddleware{sessionRepository: repository.NewSessionRepository(store)}
}

func (m *AuthMiddleware) WithCheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if !strings.HasPrefix(token, tokenPrefix) {
			warnLog(c, errors.Errorf(code.CodeUnauthorized, "Bearer token is required."))
			unauthorized(c)
			return
		}
		token = strings.TrimPrefix(token, tokenPrefix)

		userID, ok := m.sessionRepository.Get(c, token)
		if !ok {
			warnLog(c, errors.Errorf(code.CodeUnauthorized, "Request token is invalid."))
			unauthorized(c)
			return
		}
		c = (util.SetUserInfo(c, token, userID)).(*gin.Context)

		c.Next()
	}
}

func warnLog(ctx context.Context, err error) {
	logger.L.Warn(fmt.Sprint(err))
}

func unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer error=\"invalid_token\"")
	c.JSON(http.StatusUnauthorized, handler.HTTPErrUnauthorized)
	c.Abort() // ginで処理中断する際のお作法
}
