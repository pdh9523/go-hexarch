package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	"github.com/pdh9523/go-hexarch/shared/security"
	securityConfig "github.com/pdh9523/go-hexarch/shared/security/config"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	EnableCORS      bool
	EnableAuth      bool
	EnableRateLimit bool
	EnableLogging   bool
	EnableSecurity  bool
	AllowedOrigins  []string
	AllowedMethods  []string
	RateLimit       int
	MaxRequestSize  int64
}

func DefaultConfig() *Config {
	return &Config{
		EnableCORS:     true,
		EnableAuth:     true,
		EnableLogging:  true,
		EnableSecurity: true,
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		RateLimit:      10,
		MaxRequestSize: 1 << 20,
	}
}

func SetupMiddleware(
	router *gin.Engine,
	config *Config,
	responder *response.Responder,
	securityManager security.Manager,
	logger *zap.Logger,
) {
	if config == nil {
		config = DefaultConfig()
	}

	if responder == nil {
		responder = response.NewResponder()
	}

	if securityManager == nil {
		securityManager = security.NewManager(securityConfig.DefaultConfig())
	}

	if config.EnableLogging {
		router.Use(CustomLoggingMiddleware(logger))
		router.Use(RequestIDMiddleware())
	}

	if config.EnableCORS {
		if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
			router.Use(CORSMiddleware())
		} else {
			router.Use(StrictCORSMiddleware(config.AllowedOrigins))
		}
	}

	if config.EnableRateLimit {
		router.Use(IPBasedRateLimitMiddleware(responder, config.RateLimit))
	}

	router.Use(RecoveryMiddleware(responder))
}

func SetupAuthenticatedRoutes(router *gin.Engine, config *Config, responder *response.Responder, securityManager security.Manager) *gin.RouterGroup {
	if config == nil {
		config = DefaultConfig()
	}
	if responder == nil {
		responder = response.NewResponder()
	}
	if securityManager == nil {
		securityManager = security.NewManager(securityConfig.DefaultConfig())
	}

	authGroup := router.Group("/api")
	if config.EnableAuth {
		authGroup.Use(AuthMiddleware(responder, securityManager))
	}
	return authGroup
}

func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(200, gin.H{
				"status":    "ok",
				"timestamp": time.Now().Format(time.RFC3339),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
