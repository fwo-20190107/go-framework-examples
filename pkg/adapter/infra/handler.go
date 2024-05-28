package infra

import (
	"github.com/gin-gonic/gin"
)

type Handler func(c *gin.Context) *HandleError
