package jwt

import (
	"github.com/pdh9523/go-hexarch/shared/security/config"
	"testing"
	"time"
)

func TestJWTTokenManager_GenerateTokens(t *testing.T) {
	t.Run("should generate tokens successfully", func(t *testing.T) {
		// Given
		conf := config.TokenConfig{
			AccessTokenSecret:  "access-secret-key",
			RefreshTokenSecret: "refresh-secret-key",
			AccessTokenTTL:     15 * time.Minute,
			RefreshTokenTTL:    24 * time.Hour,
		}
		tm := NewTokenManager(conf)
		username := "testuser"
		role := "user"

		// When
		token, err := tm.GenerateTokens(username, role)

		// Then
		if err != nil {
			t.Fatalf("Failed to generate access token: %v", err)
		}
		if token == nil {
			t.Error("token is nil")
		}
	})
}

func TestJWTTokenManager_ValidateTokens(t *testing.T) {
	t.Run("should validate access token successfully", func(t *testing.T) {
		// Given
		conf := config.TokenConfig{
			AccessTokenSecret:  "access-secret-key",
			RefreshTokenSecret: "refresh-secret-key",
			AccessTokenTTL:     15 * time.Minute,
			RefreshTokenTTL:    24 * time.Hour,
		}
		tm := NewTokenManager(conf)
		username := "testuser"
		role := "user"
		token, _ := tm.GenerateTokens(username, role)

		// When
		accessClaims, err := tm.ValidateAccessToken(token.AccessToken)

		// Then
		if err != nil {
			t.Fatalf("Failed to validate access token: %v", err)
		}
		if accessClaims.Username != username {
			t.Errorf("Expected username %s, got %s", username, accessClaims.Username)
		}
		if accessClaims.Role != role {
			t.Errorf("Expected role %s, got %s", role, accessClaims.Role)
		}
		if accessClaims.TokenType != "access_token" {
			t.Errorf("Expected token type access, got %s", accessClaims.TokenType)
		}
	})

	t.Run("should validate refresh token successfully", func(t *testing.T) {
		// Given
		conf := config.TokenConfig{
			AccessTokenSecret:  "access-secret-key",
			RefreshTokenSecret: "refresh-secret-key",
			AccessTokenTTL:     15 * time.Minute,
			RefreshTokenTTL:    24 * time.Hour,
		}
		tm := NewTokenManager(conf)
		username := "testuser"
		role := "user"
		token, _ := tm.GenerateTokens(username, role)

		// When
		refreshClaims, err := tm.ValidateRefreshToken(token.RefreshToken)

		// Then
		if err != nil {
			t.Fatalf("Failed to validate refresh token: %v", err)
		}
		if refreshClaims.Username != username {
			t.Errorf("Expected username %s, got %s", username, refreshClaims.Username)
		}
		if refreshClaims.Role != role {
			t.Errorf("Expected role %s, got %s", role, refreshClaims.Role)
		}
		if refreshClaims.TokenType != "refresh_token" {
			t.Errorf("Expected token type refresh, got %s", refreshClaims.TokenType)
		}
	})
}

func TestJWTTokenManager_RefreshTokens(t *testing.T) {
	t.Run("should refresh tokens successfully", func(t *testing.T) {
		// Given
		conf := config.TokenConfig{
			AccessTokenSecret:  "access-secret-key",
			RefreshTokenSecret: "refresh-secret-key",
			AccessTokenTTL:     15 * time.Minute,
			RefreshTokenTTL:    24 * time.Hour,
		}
		tm := NewTokenManager(conf)
		username := "testuser"
		role := "user"
		refreshToken, _ := tm.GenerateTokens(username, role)

		// When
		newToken, err := tm.RefreshTokens(refreshToken.RefreshToken)

		// Then
		if err != nil {
			t.Fatalf("Failed to refresh tokens: %v", err)
		}
		if newToken == nil {
			t.Error("New token should not be empty")
		}
	})

	t.Run("should validate refreshed tokens", func(t *testing.T) {
		// Given
		conf := config.TokenConfig{
			AccessTokenSecret:  "access-secret-key",
			RefreshTokenSecret: "refresh-secret-key",
			AccessTokenTTL:     15 * time.Minute,
			RefreshTokenTTL:    24 * time.Hour,
		}
		tm := NewTokenManager(conf)
		username := "testuser"
		role := "user"
		token, _ := tm.GenerateTokens(username, role)

		newToken, _ := tm.RefreshTokens(token.RefreshToken)

		// When
		accessClaims, err := tm.ValidateAccessToken(newToken.AccessToken)

		// Then
		if err != nil {
			t.Fatalf("Failed to validate new access token: %v", err)
		}
		if accessClaims.Username != username {
			t.Errorf("Expected user ID %s, got %s", username, accessClaims.Username)
		}
	})
}

func TestJWTTokenManager_InvalidToken(t *testing.T) {
	t.Run("should fail with invalid token", func(t *testing.T) {
		// Given
		conf := config.TokenConfig{
			AccessTokenSecret:  "access-secret-key",
			RefreshTokenSecret: "refresh-secret-key",
			AccessTokenTTL:     15 * time.Minute,
			RefreshTokenTTL:    24 * time.Hour,
		}
		tm := NewTokenManager(conf)

		// When
		_, err := tm.ValidateAccessToken("invalid.token.string")

		// Then
		if err == nil {
			t.Error("Should fail with invalid token")
		}
	})

	t.Run("should fail with empty token", func(t *testing.T) {
		// Given
		conf := config.TokenConfig{
			AccessTokenSecret:  "access-secret-key",
			RefreshTokenSecret: "refresh-secret-key",
			AccessTokenTTL:     15 * time.Minute,
			RefreshTokenTTL:    24 * time.Hour,
		}
		tm := NewTokenManager(conf)

		// When
		_, err := tm.ValidateAccessToken("")

		// Then
		if err == nil {
			t.Error("Should fail with empty token")
		}
	})
}

func TestExtractTokenFromBearer(t *testing.T) {
	t.Run("should extract token from valid bearer token", func(t *testing.T) {
		// Given
		bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
		expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
		tm := NewTokenManager(config.TokenConfig{})
		// When
		token, err := tm.ExtractTokenFromBearer(bearerToken)

		// Then
		if err != nil {
			t.Fatalf("Failed to extract token: %v", err)
		}
		if token != expected {
			t.Errorf("Expected %s, got %s", expected, token)
		}
	})

	t.Run("should fail with invalid bearer token format", func(t *testing.T) {
		// Given
		invalidToken := "Invalid token format"
		tm := NewTokenManager(config.TokenConfig{})
		// When
		_, err := tm.ExtractTokenFromBearer(invalidToken)

		// Then
		if err == nil {
			t.Error("Should fail with invalid bearer token format")
		}
	})

	t.Run("should fail with short bearer token", func(t *testing.T) {
		// Given
		shortToken := "Bearer"
		tm := NewTokenManager(config.TokenConfig{})
		// When
		_, err := tm.ExtractTokenFromBearer(shortToken)

		// Then
		if err == nil {
			t.Error("Should fail with short bearer token")
		}
	})
}
