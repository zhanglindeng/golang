package middleware

import (
	"github.com/didip/tollbooth/config"
	"github.com/gin-gonic/gin"
	"github.com/didip/tollbooth"
	"strconv"
)


func LimitHandler(limiter *config.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter.IPLookups = []string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}
		httpError := tollbooth.LimitByRequest(limiter, c.Request)
		if httpError != nil {
			c.String(httpError.StatusCode, httpError.Message)
			c.Abort()
		} else {
			c.Writer.Header().Add("X-Rate-Limit-Limit", strconv.FormatInt(limiter.Max, 10))
			c.Writer.Header().Add("X-Rate-Limit-Duration", limiter.TTL.String())
			c.Next()
		}
	}
}

