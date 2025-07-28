package mapper

import (
	userError "github.com/pdh9523/go-hexarch/domains/user/domain/error_code"
	"net/http"
)

type ErrorCode struct {
	Code string
	err  error
}

type ErrorMapping struct {
	Code       string
	HttpStatus int
}

var errorMappings = map[error]ErrorMapping{
	userError.ErrUserNotFound: {
		Code:       "user_not_found",
		HttpStatus: http.StatusNotFound,
	},

	userError.ErrInvalidCredentials: {
		Code:       "invalid_credentials",
		HttpStatus: http.StatusUnauthorized,
	},
	userError.ErrUnauthorized: {
		Code:       "unauthorized",
		HttpStatus: http.StatusUnauthorized,
	},

	userError.ErrUserAlreadyExists: {
		Code:       "user_already_exists",
		HttpStatus: http.StatusConflict,
	},
	userError.ErrNicknameAlreadyExists: {
		Code:       "nickname_already_exists",
		HttpStatus: http.StatusConflict,
	},

	userError.ErrNicknameTooLong: {
		Code:       "nickname_too_long",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrNicknameTooShort: {
		Code:       "nickname_too_short",
		HttpStatus: http.StatusBadRequest,
	},
	userError.ErrNicknameInvalidCharacters: {
		Code:       "nickname_invalid_characters",
		HttpStatus: http.StatusBadRequest,
	},
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
}

func GetErrorMapping(err error) ErrorMapping {
	mapping, exists := errorMappings[err]
	if !exists {
		mapping = ErrorMapping{
			Code:       "unknown_error",
			HttpStatus: http.StatusInternalServerError,
		}
	}
	return mapping
}

func GetErrorCode(err error) string {
	mapping, _ := errorMappings[err]
	return mapping.Code
}
