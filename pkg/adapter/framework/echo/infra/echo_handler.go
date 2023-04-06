package infra

import (
	"context"
	"examples/pkg/adapter/infra"

	"github.com/labstack/echo/v4"
)

type EchoHandler func(ctx context.Context, c echo.Context) *infra.HandleError
