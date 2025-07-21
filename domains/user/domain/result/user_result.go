package result

import "github.com/pdh9523/go-hexarch/shared/security/jwt"

type CheckNicknameResult struct {
	IsAvailable bool
}

type CreateUserWithTokenResult struct {
	Token jwt.Token
}
