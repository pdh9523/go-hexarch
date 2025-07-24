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
	"regexp"
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

	if err := s.isNicknameExists(ctx, query.Nickname); err != nil {
		return nil, err
	}

	return result.NewCheckNicknameResult(), nil
}

func (s *UserService) SignUp(
	ctx context.Context,
	command command.SignUpCommand,
) (*result.TokenResult, error) {
	if err := s.isNicknameExists(ctx, command.Nickname); err != nil {
		return nil, err
	}

	if err := s.isUsernameExists(ctx, command.Username); err != nil {
		return nil, err
	}

	if err := validatePassword(command.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := s.securityManager.HashPassword(command.Password)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(command.Nickname, command.Username, hashedPassword)
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

	return result.NewTokenResult(token), nil
}

func (s *UserService) isUsernameExists(ctx context.Context, username string) error {
	exists, err := s.userRepository.ExistsByUsername(ctx, username)
	if err != nil {
		return err
	} else if exists {
		return errorCode.ErrUsernameAlreadyExists
	}

	return nil
}

func (s *UserService) isNicknameExists(ctx context.Context, nickname string) error {
	exists, err := s.userRepository.ExistsByNickname(ctx, nickname)
	if err != nil {
		return err
	} else if exists {
		return errorCode.ErrNicknameAlreadyExists
	}
	return nil
}

func validatePassword(password string) error {
	passwordRegex := regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]+$`)

	if len(password) < 8 {
		return errorCode.ErrPasswordTooShort
	}
	if len(password) > 32 {
		return errorCode.ErrPasswordTooLong
	}

	if !passwordRegex.MatchString(password) {
		return errorCode.ErrPasswordInvalidCharacters
	}

	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return errorCode.ErrPasswordMissingUppercase
	}

	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return errorCode.ErrPasswordMissingLowercase
	}

	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return errorCode.ErrPasswordMissingNumber
	}

	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`, password); !matched {
		return errorCode.ErrPasswordMissingSpecialChar
	}

	return nil
}
