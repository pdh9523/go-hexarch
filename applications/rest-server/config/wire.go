//go:build wireinject
// +build wireinject

package config

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/pdh9523/go-hexarch/applications/rest-server/container"
	userService "github.com/pdh9523/go-hexarch/domains/user/application/service"
	cacheAdapter "github.com/pdh9523/go-hexarch/infrastructure/cache-redis/adapter"
	cacheConfig "github.com/pdh9523/go-hexarch/infrastructure/cache-redis/config"
	gormAdapter "github.com/pdh9523/go-hexarch/infrastructure/gorm/adapter"
	gormConfig "github.com/pdh9523/go-hexarch/infrastructure/gorm/config"
	userHandler "github.com/pdh9523/go-hexarch/presentation/rest/handler/user"
	"github.com/pdh9523/go-hexarch/shared/security"
	securityConfig "github.com/pdh9523/go-hexarch/shared/security/config"
)

// DatabaseSet provides database-related dependencies
var DatabaseSet = wire.NewSet(
	gormConfig.NewGormConfig,
	ProvideGormDB,
	gormAdapter.NewUserRepositoryAdapter,
)

// CacheSet provides cache-related dependencies
var CacheSet = wire.NewSet(
	cacheConfig.NewRedisConfig,
	cacheConfig.NewRedisClient,
	cacheAdapter.NewUserCacheAdapter,
)

// SecuritySet provides security-related dependencies
var SecuritySet = wire.NewSet(
	securityConfig.DefaultConfig,
	security.NewManager,
)

// ApplicationSet provides application service dependencies
var ApplicationSet = wire.NewSet(
	userService.NewUserService,
)

// PresentationSet provides presentation layer dependencies
var PresentationSet = wire.NewSet(
	userHandler.NewUserHandler,
)

// NewContainer creates a new application container with all dependencies
func NewContainer() (container.Container, error) {
	wire.Build(
		DatabaseSet,
		CacheSet,
		SecuritySet,
		ApplicationSet,
		PresentationSet,
		ProvideContainer,
	)
	return nil, nil
}

// ProvideContainer creates a container implementation
func ProvideContainer(userHandler *userHandler.UserHandler) container.Container {
	return &containerImpl{
		userHandler: userHandler,
	}
}

// containerImpl implements the Container interface
type containerImpl struct {
	userHandler *userHandler.UserHandler
}

func (c *containerImpl) UserHandler() *userHandler.UserHandler {
	return c.userHandler
}

// ProvideGormDB is a wrapper for GormConfig.NewGormDB method
func ProvideGormDB(config *gormConfig.GormConfig) (*gorm.DB, error) {
	return config.NewGormDB()
}
