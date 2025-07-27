package adapter

import (
	"context"
	"fmt"
	userPort "github.com/pdh9523/go-hexarch/domains/user/application/port/out"
	"github.com/pdh9523/go-hexarch/infrastructure/cache-redis/client"
	"github.com/pdh9523/go-hexarch/infrastructure/cache-redis/config"
	"time"
)

type UserCacheAdapter struct {
	clientWrapper *client.RedisClientWrapper
}

const (
	RefreshTokenKeyPattern = "refresh_token:%s"
	LoginAttemptKeyPattern = "login_attempt:%s"
)

func NewUserCacheAdapter(redisClient *config.RedisClient) userPort.UserCachePort {
	return &UserCacheAdapter{
		clientWrapper: client.NewRedisClientWrapper(redisClient),
	}
}

func (a *UserCacheAdapter) SetRefreshToken(ctx context.Context, userID string, token string, ttl time.Duration) error {
	key := buildRefreshTokenKey(userID)
	return a.clientWrapper.SetDataWithTTL(ctx, key, token, ttl)
}

func (a *UserCacheAdapter) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	key := buildRefreshTokenKey(userID)
	return a.clientWrapper.Get(ctx, key)
}

func (a *UserCacheAdapter) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := buildRefreshTokenKey(userID)
	return a.clientWrapper.Delete(ctx, key)
}

func buildRefreshTokenKey(userID string) string {
	return fmt.Sprintf(RefreshTokenKeyPattern, userID)
}
