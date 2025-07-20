package in

import (
	"context"
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/query"
	"github.com/pdh9523/go-hexarch/domains/user/domain/result"
)

type UserUseCase interface {
	CheckNicknameAvailability(ext context.Context, query query.CheckNicknameQuery) (*result.CheckNicknameResult, error)
}
