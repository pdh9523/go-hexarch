package service

import (
	"context"
	usecase "github.com/pdh9523/go-hexarch/domains/user/application/port/in"
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/query"
	port "github.com/pdh9523/go-hexarch/domains/user/application/port/out"
	"github.com/pdh9523/go-hexarch/domains/user/domain"
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
	nickname, err := domain.NewNickname(query.Nickname)
	if err != nil {
		return &result.CheckNicknameResult{
			Nickname: query.Nickname,
			Status:   result.NicknamePolicyViolated,
			Reason:   err.Error(),
		}, nil
	}

	exists, err := s.userRepository.ExistsByNickname(ctx, nickname.Value())
	if err != nil {
		return nil, err
	}

	if exists {
		return &result.CheckNicknameResult{
			Nickname: nickname.Value(),
			Status:   result.NicknameDuplicated,
			Reason:   "nickname already exists",
		}, nil
	}

	return &result.CheckNicknameResult{
		Nickname: nickname.Value(),
		Status:   result.NicknameAvailable,
		Reason:   "",
	}, nil
}
