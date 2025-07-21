package domain

import (
	"github.com/pdh9523/go-hexarch/domains/user/domain/errorCode"
	"regexp"
)

type Username struct {
	value string
}

const (
	MinUsernameLength = 2
	MaxUsernameLength = 16
)

var usernameRegex = regexp.MustCompile(`^[a-z0-9]+$`)

func NewUsername(value string) (*Username, error) {
	if length := len(value); length > MaxUsernameLength {
		return nil, errorCode.ErrUsernameTooLong
	} else if length < MinUsernameLength {
		return nil, errorCode.ErrUsernameTooShort
	}

	if !usernameRegex.MatchString(value) {
		return nil, errorCode.ErrUsernameInvalidCharacters
	}

	return &Username{value: value}, nil
}

func (u *Username) ToString() string {
	return u.value
}

func (u *Username) Equals(o *Username) bool {
	if o == nil {
		return false
	}
	return u.value == o.value
}
