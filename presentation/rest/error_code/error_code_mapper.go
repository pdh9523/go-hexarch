package error_code

import (
	"errors"
	"github.com/go-playground/validator/v10"
	userError "github.com/pdh9523/go-hexarch/domains/user/domain/error_code"
	httpError "github.com/pdh9523/go-hexarch/shared/common/error_code"
	securityError "github.com/pdh9523/go-hexarch/shared/security/error_code"
	"net/http"
)

type ErrorMapping struct {
	Code       string
	HttpStatus int
}

var errorMappings = map[error]ErrorMapping{
	// === User Domain Errors ===
	// User Business Logic Errors
	userError.ErrUserNotFound: {
		Code:       "user_not_found",
		HttpStatus: http.StatusNotFound,
	},
	userError.ErrUserAlreadyExists: {
		Code:       "user_already_exists",
		HttpStatus: http.StatusConflict,
	},
	userError.ErrInvalidCredentials: {
		Code:       "invalid_credentials",
		HttpStatus: http.StatusUnauthorized,
	},
	userError.ErrUnauthorized: {
		Code:       "unauthorized",
		HttpStatus: http.StatusUnauthorized,
	},

	// Username Validation Errors
	userError.ErrUsernameAlreadyExists: {
		Code:       "username_already_exists",
		HttpStatus: http.StatusConflict,
	},
	userError.ErrUsernameTooShort: {
		Code:       "username_too_short",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrUsernameTooLong: {
		Code:       "username_too_long",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrUsernameInvalidCharacters: {
		Code:       "username_invalid_characters",
		HttpStatus: http.StatusBadRequest,
	},

	// Nickname Validation Errors
	userError.ErrNicknameAlreadyExists: {
		Code:       "nickname_already_exists",
		HttpStatus: http.StatusConflict,
	},
	userError.ErrNicknameTooShort: {
		Code:       "nickname_too_short",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrNicknameTooLong: {
		Code:       "nickname_too_long",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrNicknameInvalidCharacters: {
		Code:       "nickname_invalid_characters",
		HttpStatus: http.StatusBadRequest,
	},

	// Password Validation Errors
	userError.ErrPasswordTooShort: {
		Code:       "password_too_short",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrPasswordTooLong: {
		Code:       "password_too_long",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrPasswordInvalidCharacters: {
		Code:       "password_invalid_characters",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrPasswordMissingUppercase: {
		Code:       "password_missing_uppercase",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrPasswordMissingLowercase: {
		Code:       "password_missing_lowercase",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrPasswordMissingNumber: {
		Code:       "password_missing_number",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrPasswordMissingSpecialChar: {
		Code:       "password_missing_special_char",
		HttpStatus: http.StatusBadRequest,
	},

	// === Security/Token Errors ===
	// Token Errors (401)
	securityError.ErrTokenEmpty: {
		Code:       "token_empty",
		HttpStatus: http.StatusUnauthorized,
	},
	securityError.ErrTokenInvalid: {
		Code:       "token_invalid",
		HttpStatus: http.StatusUnauthorized,
	},
	securityError.ErrTokenExpired: {
		Code:       "token_expired",
		HttpStatus: http.StatusUnauthorized,
	},
	securityError.ErrTokenUnexpectedSigningMethod: {
		Code:       "token_invalid_signing_method",
		HttpStatus: http.StatusUnauthorized,
	},
	securityError.ErrTokenMissingAlgorithm: {
		Code:       "token_invalid_format",
		HttpStatus: http.StatusUnauthorized,
	},
	securityError.ErrTokenClaimsParseFailed: {
		Code:       "token_invalid_claims",
		HttpStatus: http.StatusUnauthorized,
	},
	securityError.ErrTokenGenerationFailed: {
		Code:       "token_generation_failed",
		HttpStatus: http.StatusInternalServerError,
	},

	// Hash/Crypto Errors (500)
	securityError.ErrInvalidHashFormat: {
		Code:       "invalid_hash_format",
		HttpStatus: http.StatusInternalServerError,
	},
	securityError.ErrIncompatibleHashVersion: {
		Code:       "incompatible_hash_version",
		HttpStatus: http.StatusInternalServerError,
	},
	securityError.ErrHashDecodingFailed: {
		Code:       "hash_decoding_failed",
		HttpStatus: http.StatusInternalServerError,
	},
	securityError.ErrPasswordMismatch: {
		Code:       "password_mismatch",
		HttpStatus: http.StatusUnauthorized,
	},

	// === HTTP Common Errors ===
	// Server Errors (5xx)
	httpError.ErrInternalServerError: {
		Code:       "internal_server_error",
		HttpStatus: http.StatusInternalServerError,
	},

	// Client Errors - Request (4xx)
	httpError.ErrInvalidRequest: {
		Code:       "invalid_request",
		HttpStatus: http.StatusBadRequest,
	},
	httpError.ErrRequestEntityTooLarge: {
		Code:       "request_entity_too_large",
		HttpStatus: http.StatusRequestEntityTooLarge,
	},
	httpError.ErrUnsupportedMediaType: {
		Code:       "unsupported_media_type",
		HttpStatus: http.StatusUnsupportedMediaType,
	},

	// Client Errors - Auth (4xx)
	httpError.ErrAuthorizationEmpty: {
		Code:       "authorization_header_empty",
		HttpStatus: http.StatusUnauthorized,
	},
	httpError.ErrAuthorizationFormat: {
		Code:       "authorization_format_error",
		HttpStatus: http.StatusUnauthorized,
	},

	// Client Errors - Rate Limiting (4xx)
	httpError.ErrRateLimitExceeded: {
		Code:       "rate_limit_exceeded",
		HttpStatus: http.StatusTooManyRequests,
	},
}

func GetErrorMapping(err error) ErrorMapping {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return handleValidationErrors(validationErrors)
	}

	mapping, exists := errorMappings[err]
	if !exists {
		mapping = ErrorMapping{
			Code:       "unknown_error",
			HttpStatus: http.StatusInternalServerError,
		}
	}
	return mapping
}

func handleValidationErrors(errs validator.ValidationErrors) ErrorMapping {
	return ErrorMapping{
		Code:       "validation_failed",
		HttpStatus: http.StatusBadRequest,
	}
}
