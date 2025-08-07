package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/applications/rest-server/config"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	"github.com/pdh9523/go-hexarch/presentation/rest/middleware"
	"github.com/pdh9523/go-hexarch/shared/security"
)

// SetupGlobalMiddlewares configures all global middlewares for the application
func SetupGlobalMiddlewares(router *gin.Engine, cfg *config.ServerConfig, securityManager security.Manager) *response.Responder {
	responder := response.NewResponder()

	// Basic Gin middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Custom middlewares (only use existing ones)
	setupCORSMiddlewares(router)
	setupSecurityMiddlewares(router, responder)

	return responder
}

// setupSecurityMiddlewares configures security-related middlewares
func setupSecurityMiddlewares(router *gin.Engine, responder *response.Responder) {
	// Add basic security headers
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})
}

// setupCORSMiddlewares configures CORS middlewares
func setupCORSMiddlewares(router *gin.Engine) {
	// Basic CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// SetupRouteSpecificMiddlewares sets up middlewares for specific route groups
func SetupRouteSpecificMiddlewares(rg *gin.RouterGroup, responder *response.Responder, securityManager security.Manager) {
	// This can be called for specific route groups that need additional middleware
	rg.Use(middleware.AuthMiddleware(responder, securityManager))
}
