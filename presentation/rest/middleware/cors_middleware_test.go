package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given regular request when processing then should set CORS headers and continue", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := CORSMiddleware()
		middleware(c)

		assert.False(t, c.IsAborted())
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "POST, OPTIONS, GET, PUT, DELETE, PATCH", w.Header().Get("Access-Control-Allow-Methods"))
	})

	t.Run("Given OPTIONS request when processing then should abort with 204", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("OPTIONS", "/test", nil)

		middleware := CORSMiddleware()
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "POST, OPTIONS, GET, PUT, DELETE, PATCH", w.Header().Get("Access-Control-Allow-Methods"))
	})

	t.Run("Given different HTTP methods when processing then should handle correctly", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

		for _, method := range methods {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(method, "/test", nil)

			middleware := CORSMiddleware()
			middleware(c)

			assert.False(t, c.IsAborted(), "Should not abort for %s method", method)
			assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		}
	})
}

func TestStrictCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	allowedOrigins := []string{
		"https://example.com",
		"https://app.example.com",
		"http://localhost:3000",
	}

	t.Run("Given allowed origin when processing then should set origin-specific CORS headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Origin", "https://example.com")

		middleware := StrictCORSMiddleware(allowedOrigins)
		middleware(c)

		assert.False(t, c.IsAborted())
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Equal(t, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "POST, OPTIONS, GET, PUT, DELETE, PATCH", w.Header().Get("Access-Control-Allow-Methods"))
	})

	t.Run("Given disallowed origin when processing then should not set origin header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Origin", "https://malicious.com")

		middleware := StrictCORSMiddleware(allowedOrigins)
		middleware(c)

		assert.False(t, c.IsAborted())
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Equal(t, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "POST, OPTIONS, GET, PUT, DELETE, PATCH", w.Header().Get("Access-Control-Allow-Methods"))
	})

	t.Run("Given no origin header when processing then should not set origin header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := StrictCORSMiddleware(allowedOrigins)
		middleware(c)

		assert.False(t, c.IsAborted())
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("Given OPTIONS request with allowed origin when processing then should abort with 204", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("OPTIONS", "/test", nil)
		c.Request.Header.Set("Origin", "http://localhost:3000")

		middleware := StrictCORSMiddleware(allowedOrigins)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("Given OPTIONS request with disallowed origin when processing then should abort with 204 without origin", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("OPTIONS", "/test", nil)
		c.Request.Header.Set("Origin", "https://evil.com")

		middleware := StrictCORSMiddleware(allowedOrigins)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("Given multiple allowed origins when checking different origins then should handle correctly", func(t *testing.T) {
		testCases := []struct {
			origin   string
			expected bool
		}{
			{"https://example.com", true},
			{"https://app.example.com", true},
			{"http://localhost:3000", true},
			{"https://notallowed.com", false},
			{"http://example.com", false},       // Different protocol
			{"https://example.com:8080", false}, // Different port
		}

		for _, tc := range testCases {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			c.Request.Header.Set("Origin", tc.origin)

			middleware := StrictCORSMiddleware(allowedOrigins)
			middleware(c)

			if tc.expected {
				assert.Equal(t, tc.origin, w.Header().Get("Access-Control-Allow-Origin"), "Origin %s should be allowed", tc.origin)
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"), "Origin %s should not be allowed", tc.origin)
			}
		}
	})

	t.Run("Given empty allowed origins list when processing then should not set origin header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Origin", "https://example.com")

		middleware := StrictCORSMiddleware([]string{})
		middleware(c)

		assert.False(t, c.IsAborted())
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})
}
