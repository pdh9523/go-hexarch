package domain

import (
	"errors"
	"regexp"
	"strings"
)

type Nickname string

const (
	MinNicknameLength = 2
	MaxNicknameLength = 20
)

var nicknameRegex = regexp.MustCompile(`^[a-zA-Z가-힣0-9]+$`)

var (
	ErrNicknameTooShort          = errors.New("nickname too short")
	ErrNicknameTooLong           = errors.New("nickname too long")
	ErrNicknameInvalidCharacters = errors.New("nickname contains invalid characters")
)

func NewNickname(value string) (Nickname, error) {
	trimmedNickname := strings.TrimSpace(value)
	return Nickname(trimmedNickname), nil
}

func (n Nickname) Value() string {
	return string(n)
}

func (n Nickname) IsEmpty() bool {
	return n == ""
}

func (n Nickname) Equals(o Nickname) bool {
	return n == o
}
