package web

import (
	"examples/pkg/adapter/framework/gin/infra"
	"examples/pkg/errors"
	"examples/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinHandler infra.GinHandler

func (fn GinHandler) Exec(c *gin.Context) {
	c.Set("Content-Type", "application/json")
	ctx := c.Request.Context()

	if err := fn(c); err != nil {
		status := errors.HTTPStatus(err.Error)
		if err.Error != nil {
			switch {
			case status >= http.StatusInternalServerError:
				logger.L.Err(ctx, fmt.Sprint(err.Error))
			case status >= http.StatusBadRequest:
				logger.L.Warn(ctx, fmt.Sprint(err.Error))
			}
		}
		c.JSON(status, err.HTTPError)
	}
}

var _ gin.HandlerFunc = (GinHandler)(nil).Exec
