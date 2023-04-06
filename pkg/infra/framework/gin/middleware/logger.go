package middleware

import (
	"regexp"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func WithLogger(opts ...logger.Option) gin.HandlerFunc {
	var reg = regexp.MustCompile(`^/swagger/*`)
	opts = append(opts, logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
		return l.Output(gin.DefaultWriter).
			With().
			Str("remoteIp", c.Request.RemoteAddr).
			Logger()
	}), logger.WithSkipPathRegexp(reg))
	return logger.SetLogger(opts...)
}
