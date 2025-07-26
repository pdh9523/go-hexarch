package error_code

import "errors"

var (
	// 인증 관련 에러
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordMismatch   = errors.New("password mismatch")

	// 해시 관련 에러 (500 에러)
	ErrInvalidHashFormat       = errors.New("invalid hash format")
	ErrIncompatibleHashVersion = errors.New("incompatible hash version")
	ErrHashDecodingFailed      = errors.New("hash decoding failed")

	// 토큰 관련 에러
	ErrTokenGenerationFailed = errors.New("token generation failed")
	ErrTokenInvalid          = errors.New("invalid token")
	ErrTokenExpired          = errors.New("token expired")
)
