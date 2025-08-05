package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given logging middleware when processing request then should not interfere with request flow", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := LoggingMiddleware()
		middleware(c)

		assert.False(t, c.IsAborted())
	})
}

func TestCustomLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given custom logging middleware when processing request then should log HTTP request details", func(t *testing.T) {
		// Create observed logger for testing
		core, logs := observer.New(zapcore.InfoLevel)
		logger := zap.New(core)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("request_id", "test-request-id")

		middleware := CustomLoggingMiddleware(logger)
		middleware(c)

		assert.False(t, c.IsAborted())

		// Verify log entry was created
		entries := logs.All()
		assert.Len(t, entries, 1)

		entry := entries[0]
		assert.Equal(t, "HTTP Request", entry.Message)
		assert.Equal(t, zapcore.InfoLevel, entry.Level)

		// Verify log fields
		fields := entry.ContextMap()
		assert.Equal(t, "test-request-id", fields["request_id"])
		assert.Equal(t, "GET", fields["method"])
		assert.Equal(t, "/test", fields["path"])
		assert.Equal(t, int64(200), fields["status_code"])
		assert.Contains(t, fields, "timestamp")
		assert.Contains(t, fields, "latency")
		assert.Contains(t, fields, "client_ip")
		assert.Contains(t, fields, "body_size")
	})

	t.Run("Given request with query parameters when processing then should log full path", func(t *testing.T) {
		core, logs := observer.New(zapcore.InfoLevel)
		logger := zap.New(core)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test?param1=value1&param2=value2", nil)
		c.Set("request_id", "test-request-id")

		middleware := CustomLoggingMiddleware(logger)
		middleware(c)

		entries := logs.All()
		assert.Len(t, entries, 1)

		fields := entries[0].ContextMap()
		assert.Equal(t, "/test?param1=value1&param2=value2", fields["path"])
	})

	t.Run("Given request without request_id when processing then should use unknown", func(t *testing.T) {
		core, logs := observer.New(zapcore.InfoLevel)
		logger := zap.New(core)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := CustomLoggingMiddleware(logger)
		middleware(c)

		entries := logs.All()
		assert.Len(t, entries, 1)

		fields := entries[0].ContextMap()
		assert.Equal(t, "unknown", fields["request_id"])
	})

	t.Run("Given request with different HTTP methods when processing then should log correct method", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

		for _, method := range methods {
			core, logs := observer.New(zapcore.InfoLevel)
			logger := zap.New(core)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(method, "/test", nil)
			c.Set("request_id", "test-request-id")

			middleware := CustomLoggingMiddleware(logger)
			middleware(c)

			entries := logs.All()
			assert.Len(t, entries, 1)

			fields := entries[0].ContextMap()
			assert.Equal(t, method, fields["method"])
		}
	})

	t.Run("Given request with different status codes when processing then should log correct status", func(t *testing.T) {
		statusCodes := []int{200, 201, 400, 404, 500}

		for _, statusCode := range statusCodes {
			core, logs := observer.New(zapcore.InfoLevel)
			logger := zap.New(core)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			c.Set("request_id", "test-request-id")

			// Set status code
			c.Status(statusCode)

			middleware := CustomLoggingMiddleware(logger)
			middleware(c)

			entries := logs.All()
			assert.Len(t, entries, 1)

			fields := entries[0].ContextMap()
			assert.Equal(t, int64(statusCode), fields["status_code"])
		}
	})

	t.Run("Given request processing time when logging then should capture latency", func(t *testing.T) {
		core, logs := observer.New(zapcore.InfoLevel)
		logger := zap.New(core)

		router := gin.New()
		router.Use(RequestIDMiddleware())
		router.Use(CustomLoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			time.Sleep(10 * time.Millisecond)
			c.Status(200)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		entries := logs.All()
		assert.Len(t, entries, 1)

		// Check that latency field exists
		fields := entries[0].ContextMap()
		assert.Contains(t, fields, "latency")

		// The latency field should be non-zero
		latencyValue := fields["latency"]
		assert.NotNil(t, latencyValue)

		// Try to get the actual latency duration from the log entry
		for _, field := range entries[0].Context {
			if field.Key == "latency" {
				if dur, ok := field.Interface.(time.Duration); ok {
					assert.True(t, dur >= 10*time.Millisecond)
				}
			}
		}
	})
}

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given request with X-Request-ID header when processing then should use existing ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("X-Request-ID", "existing-request-id")

		middleware := RequestIDMiddleware()
		middleware(c)

		requestID, exists := c.Get("request_id")
		assert.True(t, exists)
		assert.Equal(t, "existing-request-id", requestID)
		assert.Equal(t, "existing-request-id", c.Writer.Header().Get("X-Request-ID"))
		assert.False(t, c.IsAborted())
	})

	t.Run("Given request without X-Request-ID header when processing then should generate new ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := RequestIDMiddleware()
		middleware(c)

		requestID, exists := c.Get("request_id")
		assert.True(t, exists)
		assert.NotEmpty(t, requestID)

		responseHeader := c.Writer.Header().Get("X-Request-ID")
		assert.Equal(t, requestID, responseHeader)
		assert.False(t, c.IsAborted())
	})

	t.Run("Given multiple requests when processing then should generate unique IDs", func(t *testing.T) {
		var requestIDs []string

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			middleware := RequestIDMiddleware()
			middleware(c)

			requestID, exists := c.Get("request_id")
			assert.True(t, exists)
			requestIDs = append(requestIDs, requestID.(string))
		}

		// Verify all IDs are unique
		for i := 0; i < len(requestIDs); i++ {
			for j := i + 1; j < len(requestIDs); j++ {
				assert.NotEqual(t, requestIDs[i], requestIDs[j])
			}
		}
	})
}

func TestGenerateRequestID(t *testing.T) {
	t.Run("Given generateRequestID function when called then should return valid UUID", func(t *testing.T) {
		requestID := generateRequestID()

		assert.NotEmpty(t, requestID)
		assert.Len(t, requestID, 36) // UUID format: 8-4-4-4-12
		assert.Contains(t, requestID, "-")
	})

	t.Run("Given multiple calls when generating request IDs then should return unique values", func(t *testing.T) {
		var ids []string
		for i := 0; i < 100; i++ {
			ids = append(ids, generateRequestID())
		}

		idMap := make(map[string]bool)
		for _, id := range ids {
			assert.False(t, idMap[id], "Duplicate ID found: %s", id)
			idMap[id] = true
		}
	})
}
