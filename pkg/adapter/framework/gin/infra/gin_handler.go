package infra

import (
	"examples/pkg/adapter/infra"

	"github.com/gin-gonic/gin"
)

type GinHandler func(c *gin.Context) *infra.HandleError
