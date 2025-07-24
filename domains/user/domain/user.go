package domain

import (
	"regexp"
	"unicode/utf8"

	"github.com/pdh9523/go-hexarch/domains/user/domain/errorCode"
)

type User struct {
	ID       string
	Nickname string
	Username string
	Password string
	Role     Role
}

const (
	MinUsernameLength = 2
	MaxUsernameLength = 16
	MinNicknameLength = 2
	MaxNicknameLength = 20
)

var (
	usernameRegex = regexp.MustCompile(`^[a-z0-9]+$`)
	nicknameRegex = regexp.MustCompile(`^[a-zA-Z가-힣0-9]+$`)
)

func NewUser(nickname, username string) (*User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if err := ValidateNickname(nickname); err != nil {
		return nil, err
	}

	return &User{
		ID:       "",
		Nickname: nickname,
		Username: username,
		Role:     ROLE_USER,
	}, nil
}

func validateUsername(value string) error {
	if length := len(value); length > MaxUsernameLength {
		return errorCode.ErrUsernameTooLong
	} else if length < MinUsernameLength {
		return errorCode.ErrUsernameTooShort
	}

	if !usernameRegex.MatchString(value) {
		return errorCode.ErrUsernameInvalidCharacters
	}
	return nil
}

func ValidateNickname(value string) error {
	if length := utf8.RuneCountInString(value); length > MaxNicknameLength {
		return errorCode.ErrNicknameTooLong
	} else if length < MinNicknameLength {
		return errorCode.ErrNicknameTooShort
	}

	if !nicknameRegex.MatchString(value) {
		return errorCode.ErrNicknameInvalidCharacters
	}
	return nil
}
