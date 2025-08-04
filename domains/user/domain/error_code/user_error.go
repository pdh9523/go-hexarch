package error_code

import "errors"

var (
	// 사용자명 검증 에러 (400)
	ErrUsernameAlreadyExists     = errors.New("username already exists")
	ErrUsernameTooShort          = errors.New("username too short")
	ErrUsernameTooLong           = errors.New("username too long")
	ErrUsernameInvalidCharacters = errors.New("username invalid characters")

	// 닉네임 검증 에러 (400)
	ErrNicknameAlreadyExists     = errors.New("nickname already exists")
	ErrNicknameTooShort          = errors.New("nickname too short")
	ErrNicknameTooLong           = errors.New("nickname too long")
	ErrNicknameInvalidCharacters = errors.New("nickname invalid characters")

	// 비밀번호 검증 에러 (400)
	ErrPasswordTooShort           = errors.New("password too short")
	ErrPasswordTooLong            = errors.New("password too long")
	ErrPasswordInvalidCharacters  = errors.New("password contains invalid characters")
	ErrPasswordMissingUppercase   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingLowercase   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingNumber      = errors.New("password must contain at least one number")
	ErrPasswordMissingSpecialChar = errors.New("password must contain at least one special character")

	// 사용자 비즈니스 로직 에러
	ErrUserNotFound       = errors.New("user not found")      // 404
	ErrUserAlreadyExists  = errors.New("user already exists") // 409
	ErrInvalidCredentials = errors.New("invalid credentials") // 401
	ErrUnauthorized       = errors.New("unauthorized")        // 401
)
