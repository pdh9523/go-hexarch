package errorCode

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidRequest      = errors.New("invalid request")
)
