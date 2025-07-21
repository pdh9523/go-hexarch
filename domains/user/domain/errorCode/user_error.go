package errorCode

import "errors"

var (
	ErrUsernameAlreadyExists     = errors.New("username already exists")
	ErrUsernameTooShort          = errors.New("username too short")
	ErrUsernameTooLong           = errors.New("username too long")
	ErrUsernameInvalidCharacters = errors.New("username invalid characters")

	ErrNicknameAlreadyExists     = errors.New("nickname already exists")
	ErrNicknameTooShort          = errors.New("nickname too short")
	ErrNicknameTooLong           = errors.New("nickname too long")
	ErrNicknameInvalidCharacters = errors.New("nickname invalid characters")
)
