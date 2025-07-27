package result

import (
	"github.com/pdh9523/go-hexarch/domains/user/domain"
	"github.com/pdh9523/go-hexarch/shared/security/jwt"
)

type TokenResult struct {
	Token *jwt.Token
}

func NewTokenResult(token *jwt.Token) *TokenResult {
	return &TokenResult{
		Token: token,
	}
}

type UserInfoResult struct {
	Username string
	Nickname string
}

func NewUserInfoResult(user *domain.User) *UserInfoResult {
	return &UserInfoResult{
		Username: user.Username,
		Nickname: user.Nickname,
	}
}
