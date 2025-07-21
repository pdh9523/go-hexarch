package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestStandardClaims_IsExpired(t *testing.T) {
	t.Run("Given expired token When check expiration Then should return true", func(t *testing.T) {
		// Given
		claims := &StandardClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			},
		}

		// When
		result := claims.IsExpired()

		// Then
		assert.True(t, result)
	})

	t.Run("Given valid token When check expiration Then should return false", func(t *testing.T) {
		// Given
		claims := &StandardClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}

		// When
		result := claims.IsExpired()

		// Then
		assert.False(t, result)
	})

	t.Run("Given token without expiration When check expiration Then should return false", func(t *testing.T) {
		// Given
		claims := &StandardClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: nil,
			},
		}

		// When
		result := claims.IsExpired()

		// Then
		assert.False(t, result)
	})
}

func TestStandardClaims_IsValidBefore(t *testing.T) {
	t.Run("Given token not valid yet When check validity Then should return true", func(t *testing.T) {
		// Given
		claims := &StandardClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}

		// When
		result := claims.IsValidBefore()

		// Then
		assert.True(t, result)
	})

	t.Run("Given already valid token When check validity Then should return false", func(t *testing.T) {
		// Given
		claims := &StandardClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			},
		}

		// When
		result := claims.IsValidBefore()

		// Then
		assert.False(t, result)
	})

	t.Run("Given token without NotBefore When check validity Then should return true", func(t *testing.T) {
		// Given
		claims := &StandardClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				NotBefore: nil,
			},
		}

		// When
		result := claims.IsValidBefore()

		// Then
		assert.True(t, result)
	})
}

func TestStandardClaims_Getters(t *testing.T) {
	t.Run("Given claims with user data When get username Then should return correct value", func(t *testing.T) {
		// Given
		claims := &StandardClaims{
			Username: "testuser",
			Role:     "admin",
		}

		// When
		username := claims.GetUsername()

		// Then
		assert.Equal(t, "testuser", username)
	})

	t.Run("Given claims with user data When get role Then should return correct value", func(t *testing.T) {
		// Given
		claims := &StandardClaims{
			Username: "testuser",
			Role:     "admin",
		}

		// When
		role := claims.GetRole()

		// Then
		assert.Equal(t, "admin", role)
	})
}

func TestNewAccessTokenClaims(t *testing.T) {
	t.Run("Given user info and expiration When create access token claims Then should return valid claims", func(t *testing.T) {
		// Given
		username := "testuser"
		role := "admin"
		expiration := time.Hour

		// When
		claims := NewAccessTokenClaims(username, role, expiration)

		// Then
		assert.NotNil(t, claims)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, "access_token", claims.TokenType)
		assert.Equal(t, "go-hexarch", claims.Issuer)
		assert.Equal(t, username, claims.Subject)
		assert.NotNil(t, claims.IssuedAt)
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.NotBefore)

		// 만료 시간 확인
		expectedExpiry := time.Now().Add(expiration)
		actualExpiry := claims.ExpiresAt.Time
		assert.WithinDuration(t, expectedExpiry, actualExpiry, time.Second)
	})
}

func TestNewRefreshTokenClaims(t *testing.T) {
	t.Run("Given user info and expiration When create refresh token claims Then should return valid claims", func(t *testing.T) {
		// Given
		username := "testuser"
		role := "admin"
		expiration := 24 * time.Hour

		// When
		claims := NewRefreshTokenClaims(username, role, expiration)

		// Then
		assert.NotNil(t, claims)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, "refresh_token", claims.TokenType)
		assert.Equal(t, "go-hexarch", claims.Issuer)
		assert.Equal(t, username, claims.Subject)
		assert.NotNil(t, claims.IssuedAt)
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.NotBefore)

		// 만료 시간 확인
		expectedExpiry := time.Now().Add(expiration)
		actualExpiry := claims.ExpiresAt.Time
		assert.WithinDuration(t, expectedExpiry, actualExpiry, time.Second)
	})
}
