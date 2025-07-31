package error_code

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidRequest      = errors.New("invalid request")
	ErrAuthorizationEmpty  = errors.New("authorization header is empty")
	ErrAuthorizationFormat = errors.New("authorization format error")
	ErrTokenEmpty          = errors.New("token is empty")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")
)
