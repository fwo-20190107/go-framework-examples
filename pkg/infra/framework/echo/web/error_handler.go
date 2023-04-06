package web

import (
	"examples/pkg/adapter/handler"
	"examples/pkg/errors"
	"examples/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	ctx := c.Request().Context()
	status := errors.HTTPStatus(err)
	switch {
	case status >= http.StatusInternalServerError:
		{
			logger.L.Err(ctx, err.Error())
		}
	case status >= http.StatusBadRequest:
		{
			logger.L.Warn(ctx, err.Error())
		}
	default:
		{
			logger.L.Fatal(ctx, err.Error())
		}
	}

	// ここまでレスポンスが無いと言うことは構造上考えにくいが、、、
	if !c.Response().Committed {
		c.JSON(http.StatusInternalServerError, handler.HTTPErrUnauthorized)
	}
}
