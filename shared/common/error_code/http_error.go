package error_code

import "errors"

var (
	// 5xx 서버 에러
	ErrInternalServerError = errors.New("internal server error")

	// 4xx 클라이언트 에러 - 요청 관련
	ErrInvalidRequest        = errors.New("invalid request")
	ErrRequestEntityTooLarge = errors.New("request entity too large")
	ErrUnsupportedMediaType  = errors.New("unsupported media type")

	// 4xx 클라이언트 에러 - 인증/인가 관련 (HTTP 레벨)
	ErrAuthorizationEmpty  = errors.New("authorization header is empty")
	ErrAuthorizationFormat = errors.New("authorization format error")

	// 4xx 클라이언트 에러 - 제한 관련
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)
