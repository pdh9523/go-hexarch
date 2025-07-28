package security

import (
	"github.com/pdh9523/go-hexarch/shared/security/config"
	"github.com/pdh9523/go-hexarch/shared/security/crypto"
	"github.com/pdh9523/go-hexarch/shared/security/jwt"
	"time"
)

type Manager interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error

	GenerateTokens(username, role string) (*jwt.Token, error)
	ValidateAccessToken(accessToken string) (*jwt.AccessTokenClaims, error)
	ValidateRefreshToken(refreshToken string) (*jwt.RefreshTokenClaims, error)
	RefreshTokens(refreshToken string) (*jwt.Token, error)
	ExtractTokenFromBearer(bearerToken string) (string, error)

	GetRefreshTokenTTL() time.Duration
	ExtractAccessTokenClaims(accessToken string) (*jwt.AccessTokenClaims, error)
}

type DefaultManager struct {
	passwordHasher crypto.PasswordHasher
	tokenManager   jwt.Manager
}

func NewManager(config config.Config) Manager {
	return &DefaultManager{
		passwordHasher: crypto.NewArgon2PasswordHasher(),
		tokenManager:   jwt.NewTokenManager(config.Token),
	}
}

// HashPassword 비밀번호 해싱
func (sm *DefaultManager) HashPassword(password string) (string, error) {
	return sm.passwordHasher.HashPassword(password)
}

// VerifyPassword 비밀번호 검증
func (sm *DefaultManager) VerifyPassword(hashedPassword, password string) error {
	return sm.passwordHasher.VerifyPassword(hashedPassword, password)
}

// GenerateTokens 토큰 생성
func (sm *DefaultManager) GenerateTokens(username, role string) (*jwt.Token, error) {
	return sm.tokenManager.GenerateTokens(username, role)
}

// ValidateAccessToken 액세스 토큰 검증
func (sm *DefaultManager) ValidateAccessToken(tokenString string) (*jwt.AccessTokenClaims, error) {
	return sm.tokenManager.ValidateAccessToken(tokenString)
}

// ValidateRefreshToken 리프레시 토큰 검증
func (sm *DefaultManager) ValidateRefreshToken(tokenString string) (*jwt.RefreshTokenClaims, error) {
	return sm.tokenManager.ValidateRefreshToken(tokenString)
}

// RefreshTokens 토큰 갱신
func (sm *DefaultManager) RefreshTokens(refreshToken string) (*jwt.Token, error) {
	return sm.tokenManager.RefreshTokens(refreshToken)
}

// ExtractTokenFromBearer Bearer 토큰에서 실제 토큰 추출
func (sm *DefaultManager) ExtractTokenFromBearer(bearerToken string) (string, error) {
	return sm.tokenManager.ExtractTokenFromBearer(bearerToken)
}

func (sm *DefaultManager) GetRefreshTokenTTL() time.Duration {
	return sm.tokenManager.GetRefreshTokenTTL()
}

func (sm *DefaultManager) ExtractAccessTokenClaims(accessToken string) (*jwt.AccessTokenClaims, error) {
	return sm.tokenManager.ExtractAccessTokenClaims(accessToken)
}
