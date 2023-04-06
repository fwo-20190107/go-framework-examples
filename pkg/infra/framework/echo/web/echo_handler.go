package web

import (
	"examples/pkg/adapter/framework/echo/infra"
	"examples/pkg/errors"

	"github.com/labstack/echo/v4"
)

type EchoHandler infra.EchoHandler

func (fn EchoHandler) Exec(c echo.Context) error {
	c.Set("Content-Type", "application/json")
	ctx := c.Request().Context()

	if err := fn(ctx, c); err != nil {
		c.JSON(errors.HTTPStatus(err.Error), err.HTTPError)
		return err.Error
	}
	return nil
}

var _ echo.HandlerFunc = (EchoHandler)(nil).Exec
