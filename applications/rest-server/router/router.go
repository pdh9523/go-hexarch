package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/applications/rest-server/config"
	"github.com/pdh9523/go-hexarch/applications/rest-server/container"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	"github.com/pdh9523/go-hexarch/shared/security"
)

func SetupRouter(container container.Container, cfg *config.ServerConfig, securityManager security.Manager) *gin.Engine {
	if cfg.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	responder := SetupGlobalMiddlewares(router, cfg, securityManager)

	setupAPIRoutes(router, container, responder, securityManager)

	return router
}

func setupAPIRoutes(router *gin.Engine, container container.Container, responder *response.Responder, securityManager security.Manager) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "rest-server",
		})
	})

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			setupDomainRoutes(v1, container, responder, securityManager)
		}
	}
}

func setupDomainRoutes(v1 *gin.RouterGroup, container container.Container, responder *response.Responder, securityManager security.Manager) {
	userHandler := container.UserHandler()

	SetupUserRoutes(v1, userHandler)

	setupUserRoutesInternal(v1, userHandler, responder, securityManager)
}
