package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"

	"go.uber.org/zap"
)

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string { return "" })
}

func CustomLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()

		timestamp := time.Now()
		latency := timestamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		requestID, exists := c.Get("X-Request-ID")
		if !exists {
			requestID = "unknown"
		}

		logger.Info("HTTP Request",
			zap.String("request_id", requestID.(string)),
			zap.Time("timestamp", timestamp),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.Int("status_code", statusCode),
			zap.Int("body_size", bodySize),
			zap.String("path", path),
		)
	}

}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

func generateRequestID() string {
	return uuid.New().String()
}
