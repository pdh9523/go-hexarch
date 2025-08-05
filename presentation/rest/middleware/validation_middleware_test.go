package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	"github.com/stretchr/testify/assert"
)

func TestContentTypeMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	responder := response.NewResponder()
	w := httptest.NewRecorder()
	t.Run("Given a valid content type when middleware processes then should call next", func(t *testing.T) {
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/test", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		middleware := ContentTypeMiddleware(responder, "application/json")
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("Given an invalid content type and POST method when middleware processes then should call abort", func(t *testing.T) {
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/test", nil)
		c.Request.Header.Set("Content-Type", "plain text")

		middleware := ContentTypeMiddleware(responder, "application/json")
		middleware(c)

		assert.True(t, c.IsAborted())
	})

	t.Run("Given an invalid content type with GET method when middleware processes then should call next", func(t *testing.T) {
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Content-Type", "plain text")

		middleware := ContentTypeMiddleware(responder, "application/json")
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("Given an empty content type when middleware processes then should has 'application/json' type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Content-Type", "")

		middleware := ContentTypeMiddleware(responder, "application/json")
		middleware(c)

		assert.False(t, c.IsAborted())
		assert.Equal(t, c.Request.Header.Get("Content-Type"), "")
	})
}

func TestRequestSizeMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given request within size limit when middleware processes then should call next", func(t *testing.T) {
		responder := response.NewResponder()
		maxSize := int64(1000)
		contentLength := int64(500)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := strings.Repeat("a", int(contentLength))
		c.Request = httptest.NewRequest("POST", "/test", strings.NewReader(body))
		c.Request.ContentLength = contentLength

		middleware := RequestSizeMiddleware(responder, maxSize)
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("Given request exceeding size limit when middleware processes then should abort with error", func(t *testing.T) {
		responder := response.NewResponder()
		maxSize := int64(1000)
		contentLength := int64(1500)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := strings.Repeat("a", int(contentLength))
		c.Request = httptest.NewRequest("POST", "/test", strings.NewReader(body))
		c.Request.ContentLength = contentLength

		middleware := RequestSizeMiddleware(responder, maxSize)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	})

	t.Run("Given request at exact size limit when middleware processes then should call next", func(t *testing.T) {
		responder := response.NewResponder()
		maxSize := int64(1000)
		contentLength := int64(1000)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := strings.Repeat("a", int(contentLength))
		c.Request = httptest.NewRequest("POST", "/test", strings.NewReader(body))
		c.Request.ContentLength = contentLength

		middleware := RequestSizeMiddleware(responder, maxSize)
		middleware(c)

		assert.False(t, c.IsAborted())
	})
}
