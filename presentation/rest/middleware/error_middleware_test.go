package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given panic with error when processing then should recover and respond with error", func(t *testing.T) {
		responder := response.NewResponder()

		router := gin.New()
		router.Use(RecoveryMiddleware(responder))
		router.GET("/test", func(c *gin.Context) {
			panic(errors.New("test panic error"))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("Given panic with string when processing then should recover and respond with error", func(t *testing.T) {
		responder := response.NewResponder()

		router := gin.New()
		router.Use(RecoveryMiddleware(responder))
		router.GET("/test", func(c *gin.Context) {
			panic("test panic string")
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("Given panic with other type when processing then should recover and respond with error", func(t *testing.T) {
		responder := response.NewResponder()

		router := gin.New()
		router.Use(RecoveryMiddleware(responder))
		router.GET("/test", func(c *gin.Context) {
			panic(123) // panic with integer
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("Given no panic when processing then should handle normally", func(t *testing.T) {
		responder := response.NewResponder()

		router := gin.New()
		router.Use(RecoveryMiddleware(responder))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "ok")
	})

	t.Run("Given panic in middleware chain when processing then should recover", func(t *testing.T) {
		responder := response.NewResponder()

		router := gin.New()
		router.Use(RecoveryMiddleware(responder))
		router.Use(func(c *gin.Context) {
			panic("middleware panic")
		})
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("Given panic with nil when processing then should recover and respond with error", func(t *testing.T) {
		responder := response.NewResponder()

		router := gin.New()
		router.Use(RecoveryMiddleware(responder))
		router.GET("/test", func(c *gin.Context) {
			var nilErr error
			panic(nilErr)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("Given panic with struct when processing then should recover and format error", func(t *testing.T) {
		responder := response.NewResponder()

		type customError struct {
			Code    int
			Message string
		}

		router := gin.New()
		router.Use(RecoveryMiddleware(responder))
		router.GET("/test", func(c *gin.Context) {
			panic(customError{Code: 123, Message: "custom error"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("Given multiple panics when processing then each should be handled", func(t *testing.T) {
		responder := response.NewResponder()

		router := gin.New()
		router.Use(RecoveryMiddleware(responder))

		router.GET("/panic-error", func(c *gin.Context) {
			panic(errors.New("error panic"))
		})

		router.GET("/panic-string", func(c *gin.Context) {
			panic("string panic")
		})

		router.GET("/panic-number", func(c *gin.Context) {
			panic(42)
		})

		testCases := []string{"/panic-error", "/panic-string", "/panic-number"}

		for _, path := range testCases {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code, "Path %s should return 500", path)
			assert.Contains(t, w.Body.String(), "error", "Path %s should contain error", path)
		}
	})
}
