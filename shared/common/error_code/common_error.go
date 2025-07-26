package error_code

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidRequest      = errors.New("invalid request")
)
