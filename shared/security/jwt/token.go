package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pdh9523/go-hexarch/shared/security/config"
)

type Manager interface {
	GenerateTokens(username, role string) (*Token, error)
	ValidateAccessToken(accessToken string) (*AccessTokenClaims, error)
	ValidateRefreshToken(refreshToken string) (*RefreshTokenClaims, error)
	RefreshTokens(refreshToken string) (*Token, error)
	ExtractTokenFromBearer(bearerToken string) (string, error)
	GetRefreshTokenTTL() time.Duration
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type TokenManager struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
	signingMethod      jwt.SigningMethod
}

func NewTokenManager(config config.TokenConfig) *TokenManager {
	return &TokenManager{
		accessTokenSecret:  []byte(config.AccessTokenSecret),
		refreshTokenSecret: []byte(config.RefreshTokenSecret),
		accessTokenTTL:     config.AccessTokenTTL,
		refreshTokenTTL:    config.RefreshTokenTTL,
		signingMethod:      jwt.SigningMethodHS256,
	}
}

func (t *TokenManager) GenerateTokens(username, role string) (*Token, error) {
	accessToken, err := t.generateAccessToken(username, role)
	if err != nil {
		return nil, err
	}
	refreshToken, err := t.generateRefreshToken(username, role)
	if err != nil {
		return nil, err
	}

	return &Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (t *TokenManager) ValidateAccessToken(accessToken string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method" + token.Header["alg"].(string))
		}
		return t.accessTokenSecret, nil
	})

	if err != nil {
		return nil, errors.New("failed to validate access token: " + err.Error())
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	if claims.TokenType != "access_token" {
		return nil, errors.New("not an access token")
	}
	return claims, nil
}

func (t *TokenManager) ValidateRefreshToken(refreshToken string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method" + token.Header["alg"].(string))
		}
		return t.refreshTokenSecret, nil
	})

	if err != nil {
		return nil, errors.New("failed to validate refresh token: " + err.Error())
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	if claims.TokenType != "refresh_token" {
		return nil, errors.New("not an refresh token")
	}
	return claims, nil
}

func (t *TokenManager) RefreshTokens(refreshToken string) (*Token, error) {
	claims, err := t.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token" + err.Error())
	}

	tokens, err := t.GenerateTokens(claims.Username, claims.Role)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (t *TokenManager) generateAccessToken(username, role string) (string, error) {
	claims := NewAccessTokenClaims(username, role, t.accessTokenTTL)
	token := jwt.NewWithClaims(t.signingMethod, claims)
	return token.SignedString(t.accessTokenSecret)
}

func (t *TokenManager) generateRefreshToken(username, role string) (string, error) {
	claims := NewRefreshTokenClaims(username, role, t.refreshTokenTTL)
	token := jwt.NewWithClaims(t.signingMethod, claims)
	return token.SignedString(t.refreshTokenSecret)
}

func (t *TokenManager) ExtractTokenFromBearer(bearerToken string) (string, error) {
	if len(bearerToken) < 7 || bearerToken[:7] != "Bearer " {
		return "", errors.New("invalid bearer token format")
	}
	return bearerToken[7:], nil
}

func (t *TokenManager) GetRefreshTokenTTL() time.Duration {
	return t.refreshTokenTTL
}
