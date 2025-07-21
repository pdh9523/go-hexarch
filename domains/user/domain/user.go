package domain

import (
	"errors"
	"strings"
)

type User struct {
	ID       string
	Nickname *Nickname
	Role     Role
}
