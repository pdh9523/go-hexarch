package result

import "github.com/pdh9523/go-hexarch/shared/security/jwt"

type CheckNicknameResult struct {
	IsAvailable bool
}

func NewCheckNicknameResult() *CheckNicknameResult {
	return &CheckNicknameResult{
		IsAvailable: true,
	}
}

type TokenResult struct {
	Token *jwt.Token
}

func NewTokenResult(token *jwt.Token) *TokenResult {
	return &TokenResult{
		Token: token,
	}
}
