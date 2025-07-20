package domain

import (
	"strings"
)

type Nickname string

func NewNickname(value string) (Nickname, error) {
	trimmedNickname := strings.TrimSpace(value)

	return Nickname(trimmedNickname), nil
}

func (n Nickname) Value() string {
	return string(n)
}

func (n Nickname) IsEmpty() bool {
	return n == Nickname("")
}

func (n Nickname) Equals(o Nickname) bool {
	return n == o
}
