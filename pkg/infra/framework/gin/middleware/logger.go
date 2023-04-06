package middleware

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func WithLogger(opts ...logger.Option) gin.HandlerFunc {
	opts = append(opts, logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
		return l.Output(gin.DefaultWriter).
			With().
			Str("remoteIp", c.Request.RemoteAddr).
			Logger()
	}))
	return logger.SetLogger(opts...)
}
