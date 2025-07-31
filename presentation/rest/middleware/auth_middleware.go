package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	httpError "github.com/pdh9523/go-hexarch/shared/common/error_code"
	"github.com/pdh9523/go-hexarch/shared/security"
	securityError "github.com/pdh9523/go-hexarch/shared/security/error_code"
)

func AuthMiddleware(responder *response.Responder, securityManager security.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			responder.AbortWithError(c, httpError.ErrAuthorizationEmpty)
			return
		}
		token, err := securityManager.ExtractTokenFromBearer(authHeader)
		if err != nil {
			responder.AbortWithError(c, securityError.ErrTokenInvalid)
			return
		}

		if token == "" {
			responder.AbortWithError(c, httpError.ErrTokenEmpty)
			return
		}

		accessTokenClaims, err := securityManager.ExtractAccessTokenClaims(token)
		if err != nil {
			responder.AbortWithError(c, httpError.ErrInvalidCredentials)
			return
		}

		c.Set("user_id", accessTokenClaims.GetUserID())
		c.Set("role", accessTokenClaims.GetRole())

		c.Next()
	}
}

func OptionalAuthMiddleware(securityManager security.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			token, err := securityManager.ExtractTokenFromBearer(authHeader)
			if err == nil {
				accessTokenClaims, err := securityManager.ExtractAccessTokenClaims(token)
				if err == nil {
					c.Set("user_id", accessTokenClaims.GetUserID())
					c.Set("role", accessTokenClaims.GetRole())
				}
			}
		}
		c.Next()
	}
}

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// XSS Protection
		c.Header("X-XSS-Protection", "1; mode=block")
		// Content Type Options
		c.Header("X-Content-Type-Options", "nosniff")
		// Frame Options
		c.Header("X-Frame-Options", "DENY")
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// Content Security Policy (기본적인 정책)
		c.Header("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}
