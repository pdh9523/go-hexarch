package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	commonError "github.com/pdh9523/go-hexarch/shared/common/error_code"
	"github.com/pdh9523/go-hexarch/shared/security"
	securityError "github.com/pdh9523/go-hexarch/shared/security/error_code"
)

func AuthMiddleware(responder *response.Responder, securityManager security.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			responder.AbortWithError(c, commonError.ErrAuthorizationEmpty)
			return
		}
		token, err := securityManager.ExtractTokenFromBearer(authHeader)
		if err != nil {
			responder.AbortWithError(c, securityError.ErrTokenInvalid)
		}

		if token == "" {
			responder.AbortWithError(c, commonError.ErrTokenEmpty)
			return
		}

		accessTokenClaims, err := securityManager.ExtractAccessTokenClaims(token)
		if err != nil {
			responder.AbortWithError(c, commonError.ErrInvalidCredentials)
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
