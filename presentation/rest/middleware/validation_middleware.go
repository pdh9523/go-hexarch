package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	commonError "github.com/pdh9523/go-hexarch/shared/common/error_code"
	"strings"
)

func ContentTypeMiddleware(responder *response.Responder, allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")

			if contentType == "" {
				contentType = "application/json"
			}

			allowed := false
			for _, allowedType := range allowedTypes {
				if strings.Contains(contentType, allowedType) {
					allowed = true
					break
				}
			}
			if !allowed {
				responder.AbortWithError(c, commonError.ErrUnsupportedMediaType)
				return
			}
		}
		c.Next()
	}
}

func JSONOnlyMiddleware(responder *response.Responder) gin.HandlerFunc {
	return ContentTypeMiddleware(responder, "application/json")
}

func RequestSizeMiddleware(responder *response.Responder, maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			responder.AbortWithError(c, commonError.ErrRequestEntityTooLarge)
			return
		}
		c.Next()
	}
}
