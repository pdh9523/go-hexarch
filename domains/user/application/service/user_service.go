package service

import (
	"context"
	usecase "github.com/pdh9523/go-hexarch/domains/user/application/port/in"
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/command"
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/query"
	port "github.com/pdh9523/go-hexarch/domains/user/application/port/out"
	"github.com/pdh9523/go-hexarch/domains/user/domain"
	"github.com/pdh9523/go-hexarch/domains/user/domain/errorCode"
	"github.com/pdh9523/go-hexarch/domains/user/domain/result"
	"github.com/pdh9523/go-hexarch/shared/security"
)

type UserService struct {
	securityManager security.Manager
	userRepository  port.UserRepositoryPort
}

func NewUserService(
	securityManager security.Manager,
	userRepository port.UserRepositoryPort,
) usecase.UserUseCase {
	return &UserService{
		securityManager: securityManager,
		userRepository:  userRepository,
	}
}

func (s *UserService) CheckNicknameAvailability(
	ctx context.Context,
	query query.CheckNicknameQuery,
) (*result.CheckNicknameResult, error) {
	if err := domain.ValidateNickname(query.Nickname); err != nil {
		return nil, err
	}

	exists, err := s.userRepository.ExistsByNickname(ctx, query.Nickname)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errorCode.ErrNicknameAlreadyExists
	}

	return &result.CheckNicknameResult{IsAvailable: !exists}, nil
}

func (s *UserService) SignUp(
	ctx context.Context,
	command command.SignUpCommand,
) (*result.TokenResult, error) {
	exists, err := s.userRepository.ExistsByUsername(ctx, command.Username)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, errorCode.ErrUsernameAlreadyExists
	}

	user, err := domain.NewUser(command.Nickname, command.Username)
	if err != nil {
		return nil, err
	}

	savedUser, err := s.userRepository.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := s.securityManager.GenerateTokens(savedUser.Username, savedUser.Role.ToString())
	if err != nil {
		return nil, err
	}
	return &result.TokenResult{Token: *token}, nil
}
