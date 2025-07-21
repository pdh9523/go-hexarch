package domain

import (
	"errors"
	"strings"
)

type User struct {
	Nickname *Nickname
	Username *Username
	Role     Role
}
