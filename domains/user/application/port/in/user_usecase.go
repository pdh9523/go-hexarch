package in

import (
	"context"
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/command"
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/query"
	"github.com/pdh9523/go-hexarch/domains/user/domain/result"
)

type UserUseCase interface {
	CheckNicknameAvailability(cxt context.Context, query query.CheckNicknameQuery) (*result.CheckNicknameResult, error)
	SignUp(ctx context.Context, command command.SignUpCommand) (*result.TokenResult, error)
	SignIn(ctx context.Context, command command.SignInCommand) (*result.TokenResult, error)
	ChangePassword(ctx context.Context, command command.ChangePasswordCommand) error
	ChangeNickname(ctx context.Context, command command.ChangeNicknameCommand) error
}
