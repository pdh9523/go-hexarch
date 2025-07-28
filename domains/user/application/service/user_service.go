package service

import (
	"context"
	usecase "github.com/pdh9523/go-hexarch/domains/user/application/port/in"
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/command"
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/query"
	port "github.com/pdh9523/go-hexarch/domains/user/application/port/out"
	"github.com/pdh9523/go-hexarch/domains/user/domain"
	"github.com/pdh9523/go-hexarch/domains/user/domain/error_code"
	"github.com/pdh9523/go-hexarch/domains/user/domain/result"
	"github.com/pdh9523/go-hexarch/shared/security"
	"regexp"
)

type UserService struct {
	securityManager security.Manager
	userRepository  port.UserRepositoryPort
	userCache       port.UserCachePort
}

func NewUserService(
	securityManager security.Manager,
	userRepository port.UserRepositoryPort,
	userCache port.UserCachePort,
) usecase.UserUseCase {
	return &UserService{
		securityManager: securityManager,
		userRepository:  userRepository,
		userCache:       userCache,
	}
}
func (s *UserService) GetMyInfo(
	ctx context.Context,
) (*result.UserInfoResult, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, error_code.ErrInvalidCredentials
	}
	if userID == "" {
		return nil, error_code.ErrInvalidCredentials
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return result.NewUserInfoResult(user), nil
}

func (s *UserService) CheckNicknameAvailability(
	ctx context.Context,
	query query.CheckNicknameQuery,
) error {
	if err := domain.ValidateNickname(query.Nickname); err != nil {
		return err
	}

	if err := s.isNicknameExists(ctx, query.Nickname); err != nil {
		return err
	}

	return nil
}

func (s *UserService) CheckUsernameAvailability(
	ctx context.Context,
	query query.CheckUsernameQuery,
) error {
	if err := domain.ValidateUsername(query.Username); err != nil {
		return err
	}

	if err := s.isUsernameExists(ctx, query.Username); err != nil {
		return err
	}

	return nil
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

	token, err := s.securityManager.GenerateTokens(savedUser.ID, savedUser.Role.ToString())
	if err != nil {
		return nil, err
	}

	if err := s.userCache.SetRefreshToken(ctx, user.ID, token.RefreshToken, s.securityManager.GetRefreshTokenTTL()); err != nil {
		return nil, err
	}

	return result.NewTokenResult(token), nil
}

func (s *UserService) SignIn(
	ctx context.Context,
	command command.SignInCommand,
) (*result.TokenResult, error) {
	if err := domain.ValidateUsername(command.Username); err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindByUsername(ctx, command.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, error_code.ErrInvalidCredentials
	}

	if err = s.securityManager.VerifyPassword(user.Password, command.Password); err != nil {
		return nil, error_code.ErrInvalidCredentials
	}

	token, err := s.securityManager.GenerateTokens(user.ID, user.Role.ToString())
	if err != nil {
		return nil, err
	}

	if err := s.userCache.SetRefreshToken(ctx, user.ID, token.RefreshToken, s.securityManager.GetRefreshTokenTTL()); err != nil {
		return nil, err
	}

	return result.NewTokenResult(token), nil
}

func (s *UserService) SignOut(ctx context.Context) error {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return error_code.ErrInvalidCredentials
	}

	return s.userCache.DeleteRefreshToken(ctx, userID)
}

func (s *UserService) ChangePassword(
	ctx context.Context,
	command command.ChangePasswordCommand,
) error {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return error_code.ErrUnauthorized
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return error_code.ErrUserNotFound
	}

	if err = s.securityManager.VerifyPassword(user.Password, command.CurrentPassword); err != nil {
		return error_code.ErrInvalidCredentials
	}

	if err = validatePassword(command.NewPassword); err != nil {
		return err
	}

	hashedPassword, err := s.securityManager.HashPassword(command.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	_, err = s.userRepository.Update(ctx, user)
	return err
}

func (s *UserService) ChangeNickname(
	ctx context.Context,
	command command.ChangeNicknameCommand,
) error {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return error_code.ErrUnauthorized
	}

	if err := domain.ValidateNickname(command.Nickname); err != nil {
		return err
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return error_code.ErrUserNotFound
	}

	if user.Nickname != command.Nickname {
		if err := s.isNicknameExists(ctx, command.Nickname); err != nil {
			return err
		}
	}

	user.Nickname = command.Nickname
	_, err = s.userRepository.Update(ctx, user)
	return err
}

func (s *UserService) isUsernameExists(ctx context.Context, username string) error {
	exists, err := s.userRepository.ExistsByUsername(ctx, username)
	if err != nil {
		return err
	} else if exists {
		return error_code.ErrUsernameAlreadyExists
	}

	return nil
}

func (s *UserService) isNicknameExists(ctx context.Context, nickname string) error {
	exists, err := s.userRepository.ExistsByNickname(ctx, nickname)
	if err != nil {
		return err
	} else if exists {
		return error_code.ErrNicknameAlreadyExists
	}
	return nil
}

func validatePassword(password string) error {
	passwordRegex := regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]+$`)

	if len(password) < 8 {
		return error_code.ErrPasswordTooShort
	}
	if len(password) > 32 {
		return error_code.ErrPasswordTooLong
	}

	if !passwordRegex.MatchString(password) {
		return error_code.ErrPasswordInvalidCharacters
	}

	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return error_code.ErrPasswordMissingUppercase
	}

	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return error_code.ErrPasswordMissingLowercase
	}

	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return error_code.ErrPasswordMissingNumber
	}

	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`, password); !matched {
		return error_code.ErrPasswordMissingSpecialChar
	}

	return nil
}
