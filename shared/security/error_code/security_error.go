package error_code

import "errors"

var (
	// 암호화/해시 관련 에러 (500 에러)
	ErrInvalidHashFormat       = errors.New("invalid hash format")
	ErrIncompatibleHashVersion = errors.New("incompatible hash version")
	ErrHashDecodingFailed      = errors.New("hash decoding failed")
	ErrPasswordMismatch        = errors.New("password mismatch")

	// 토큰 관련 에러 (401 에러)
	ErrTokenGenerationFailed        = errors.New("token generation failed")
	ErrTokenEmpty                   = errors.New("token is empty")
	ErrTokenInvalid                 = errors.New("invalid token")
	ErrTokenExpired                 = errors.New("token expired")
	ErrTokenUnexpectedSigningMethod = errors.New("unexpected token signing method")
	ErrTokenMissingAlgorithm        = errors.New("missing algorithm in token header")
	ErrTokenClaimsParseFailed       = errors.New("failed to parse token claims")
)
