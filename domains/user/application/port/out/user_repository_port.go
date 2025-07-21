package out

import (
	"context"
	"github.com/pdh9523/go-hexarch/domains/user/domain"
)

type UserRepositoryPort interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByNickname(ctx context.Context, nickname string) (*domain.User, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	ExistsByNickname(ctx context.Context, nickname string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}
