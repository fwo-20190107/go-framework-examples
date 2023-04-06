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

	"github.com/labstack/echo/v4"
)

const tokenPrefix = "Bearer "

var Auth *AuthMiddleware

type AuthMiddleware struct {
	sessionRepository IRepository.SessionRepository
}

func InitAuthMiddleware(store infra.LocalStore) {
	Auth = &AuthMiddleware{sessionRepository: repository.NewSessionRepository(store)}
}

func (m *AuthMiddleware) WithCheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var header struct {
			Token string `header:"Authorization"`
		}
		if err := c.Bind(&header); err != nil {
			return err
		}
		if !strings.HasPrefix(header.Token, tokenPrefix) {
			err := errors.Errorf(code.CodeUnauthorized, "Bearer token is required.")
			unauthorized(c)
			return err
		}
		header.Token = strings.TrimPrefix(header.Token, tokenPrefix)

		userID, ok := m.sessionRepository.Get(ctx, header.Token)
		if !ok {
			err := errors.Errorf(code.CodeUnauthorized, "Request token is invalid.")
			unauthorized(c)
			return err
		}
		ctx = util.SetUserInfo(ctx, header.Token, userID)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func warnLog(ctx context.Context, err error) {
	logger.L.Warn(fmt.Sprint(err))
}

func unauthorized(c echo.Context) {
	c.Response().Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
	c.JSON(http.StatusUnauthorized, handler.HTTPErrUnauthorized)
}
