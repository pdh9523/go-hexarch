package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// StandardClaims 표준 JWT 클레임
type StandardClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AccessTokenClaims 액세스 토큰을 위한 클레임
type AccessTokenClaims struct {
	StandardClaims
	TokenType string `json:"token_type"`
}

// RefreshTokenClaims 리프레시 토큰을 위한 클레임
type RefreshTokenClaims struct {
	StandardClaims
	TokenType string `json:"token_type"`
}

// NewAccessTokenClaims 새로운 액세스 토큰 클레임 생성
func NewAccessTokenClaims(userID, username, role string, expiration time.Duration) *AccessTokenClaims {
	now := time.Now()
	return &AccessTokenClaims{
		StandardClaims: StandardClaims{
			UserID:   userID,
			Username: username,
			Role:     role,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "go-hexarch",
				Subject:   userID,
			},
		},
		TokenType: "access_token",
	}
}

// NewRefreshTokenClaims 새로운 리프레시 토큰 클레임 생성
func NewRefreshTokenClaims(userID, username, role string, expiration time.Duration) *RefreshTokenClaims {
	now := time.Now()
	return &RefreshTokenClaims{
		StandardClaims: StandardClaims{
			UserID:   userID,
			Username: username,
			Role:     role,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "go-hexarch",
				Subject:   userID,
			},
		},
		TokenType: "refresh_token",
	}
}

// IsExpired 토큰이 만료되었는지 확인
func (c *StandardClaims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(c.ExpiresAt.Time)
}

// IsValidBefore 토큰이 아직 유효하지 않은지 확인
func (c *StandardClaims) IsValidBefore() bool {
	if c.NotBefore == nil {
		return true
	}
	return time.Now().Before(c.NotBefore.Time)
}

// GetUserID 사용자 ID 반환
func (c *StandardClaims) GetUserID() string {
	return c.UserID
}

// GetUsername 사용자명 반환
func (c *StandardClaims) GetUsername() string {
	return c.Username
}

// GetRole 사용자 역할 반환
func (c *StandardClaims) GetRole() string {
	return c.Role
}
