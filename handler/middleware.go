package handler

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/phpgao/tlog"
)

// GinLogger forked from github.com/toorop/gin-logrus
func GinLogger() gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknow"
	}
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		logger := tlog.With([]tlog.Field{
			{"hostname", hostname},
			{"statusCode", statusCode},
			{"latency", latency},
			{"clientIP", clientIP},
			{"method", c.Request.Method},
			{"path", path},
			{"referer", referer},
			{"dataLength", dataLength},
			{"userAgent", clientUserAgent},
		}...)

		if len(c.Errors) > 0 {
			logger.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)",
				clientIP, hostname, time.Now().Format(time.DateTime), c.Request.Method,
				path, statusCode, dataLength, referer, clientUserAgent, latency)
			if statusCode >= http.StatusInternalServerError {
				logger.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				logger.Warn(msg)
			} else {
				logger.Info(msg)
			}
		}
	}
}
