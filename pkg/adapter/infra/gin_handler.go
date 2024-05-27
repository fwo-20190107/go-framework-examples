package infra

import (
	"github.com/gin-gonic/gin"
)

type GinHandler func(c *gin.Context) *HandleError
