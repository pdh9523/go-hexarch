package domain

import (
	"github.com/pdh9523/go-hexarch/domains/user/domain/errorCode"
	"regexp"
	"unicode/utf8"
)

type Nickname struct {
	value string
}

const (
	MinNicknameLength = 2
	MaxNicknameLength = 20
)

var nicknameRegex = regexp.MustCompile(`^[a-zA-Z가-힣0-9]+$`)

func NewNickname(value string) (*Nickname, error) {
	if length := utf8.RuneCountInString(value); length > MaxNicknameLength {
		return nil, errorCode.ErrNicknameTooLong
	} else if length < MinNicknameLength {
		return nil, errorCode.ErrNicknameTooShort
	}

	if !nicknameRegex.MatchString(value) {
		return nil, errorCode.ErrNicknameInvalidCharacters
	}

	return &Nickname{value: value}, nil
}

func (n *Nickname) ToString() string {
	return n.value
}

func (n *Nickname) IsEmpty() bool {
	return n.value == ""
}

func (n *Nickname) Equals(o *Nickname) bool {
	if o == nil {
		return false
	}
	return n.value == o.value
}
