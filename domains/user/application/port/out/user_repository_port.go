package out

import (
	"context"
	"os/user"
)

type UserRepositoryPort interface {
	FindById(ctx context.Context, id string) (*user.User, error)
	FindByNickname(ctx context.Context, nickname string) (*user.User, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	ExistsByNickname(ctx context.Context, nickname string) (bool, error)
	Save(ctx context.Context, user *user.User) error
	Update(ctx context.Context, user *user.User) error
}
