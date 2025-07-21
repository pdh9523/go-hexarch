package errorCode

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserNotFound          = errors.New("user not found")

	ErrNicknameAlreadyExists     = errors.New("nickname already exists")
	ErrNicknameTooShort          = errors.New("nickname too short")
	ErrNicknameTooLong           = errors.New("nickname too long")
	ErrNicknameInvalidCharacters = errors.New("nickname invalid characters")
)
