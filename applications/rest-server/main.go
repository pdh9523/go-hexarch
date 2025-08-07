package main

import (
	"log"

	"github.com/pdh9523/go-hexarch/applications/rest-server/config"
	"github.com/pdh9523/go-hexarch/applications/rest-server/router"
	envConfig "github.com/pdh9523/go-hexarch/shared/common/config"
	"github.com/pdh9523/go-hexarch/shared/security"
	securityConfig "github.com/pdh9523/go-hexarch/shared/security/config"
)

func main() {
	// Load environment variables from .env file
	envConfig.LoadEnv()

	// Initialize the application container with all dependencies
	container, err := config.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// Server configuration
	cfg := config.NewDefaultServerConfig()

	// Security manager
	securityManager := security.NewManager(securityConfig.DefaultConfig())

	// Setup router with all dependencies
	r := router.SetupRouter(container, cfg, securityManager)

	log.Printf("Starting server on port %s in %s mode", cfg.Port, cfg.Mode)
	log.Printf("Available routes:")
	log.Printf("  GET  /health")
	log.Printf("  USER DOMAIN:")
	log.Printf("    GET  /api/v1/user/nickname/check")
	log.Printf("    POST /api/v1/user/ (signup)")
	log.Printf("    POST /api/v1/users/signup")
	log.Printf("    POST /api/v1/users/signin")
	log.Printf("    GET  /api/v1/users/username/check")
	log.Printf("    GET  /api/v1/users/nickname/check")
	log.Printf("    GET  /api/v1/users/me (protected)")
	log.Printf("    POST /api/v1/users/signout (protected)")
	log.Printf("    PUT  /api/v1/users/password (protected)")
	log.Printf("    PUT  /api/v1/users/nickname (protected)")

	// Start the server
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
