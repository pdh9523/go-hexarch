package out

import (
	"context"
	"time"
)

type UserCachePort interface {
	SetRefreshToken(ctx context.Context, userID string, token string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
}
