package error_code

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

	ErrPasswordTooShort           = errors.New("password too short")
	ErrPasswordTooLong            = errors.New("password too long")
	ErrPasswordInvalidCharacters  = errors.New("password contains invalid characters")
	ErrPasswordMissingUppercase   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingLowercase   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingNumber      = errors.New("password must contain at least one number")
	ErrPasswordMissingSpecialChar = errors.New("password must contain at least one special character")

	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
