package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pdh9523/go-hexarch/shared/security/config"
	securityError "github.com/pdh9523/go-hexarch/shared/security/error_code"
)

type Manager interface {
	GenerateTokens(userID, role string) (*Token, error)
	ValidateAccessToken(accessToken string) (*AccessTokenClaims, error)
	ValidateRefreshToken(refreshToken string) (*RefreshTokenClaims, error)
	RefreshTokens(refreshToken string) (*Token, error)
	ExtractTokenFromBearer(bearerToken string) (string, error)
	GetRefreshTokenTTL() time.Duration
	ExtractAccessTokenClaims(accessToken string) (*AccessTokenClaims, error)
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

func (t *TokenManager) GenerateTokens(userID, role string) (*Token, error) {
	accessToken, err := t.generateAccessToken(userID, role)
	if err != nil {
		return nil, err
	}
	refreshToken, err := t.generateRefreshToken(userID, role)
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

	tokens, err := t.GenerateTokens(claims.UserID, claims.Role)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (t *TokenManager) generateAccessToken(userID, role string) (string, error) {
	claims := NewAccessTokenClaims(userID, role, t.accessTokenTTL)
	token := jwt.NewWithClaims(t.signingMethod, claims)
	return token.SignedString(t.accessTokenSecret)
}

func (t *TokenManager) generateRefreshToken(userID, role string) (string, error) {
	claims := NewRefreshTokenClaims(userID, role, t.refreshTokenTTL)
	token := jwt.NewWithClaims(t.signingMethod, claims)
	return token.SignedString(t.refreshTokenSecret)
}

func (t *TokenManager) ExtractTokenFromBearer(bearerToken string) (string, error) {
	if len(bearerToken) < 7 || bearerToken[:7] != "Bearer " {
		return "", errors.New("invalid bearer token format")
	}
	return bearerToken[7:], nil
}

func (t *TokenManager) ExtractAccessTokenClaims(accessToken string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 알고리즘 검증 - 화이트리스트 방식
		if token.Method != jwt.SigningMethodHS256 {
			alg, ok := token.Header["alg"].(string)
			if !ok {
				return nil, securityError.ErrTokenMissingAlgorithm
			}
			return nil, fmt.Errorf("%w: %s", securityError.ErrTokenUnexpectedSigningMethod, alg)
		}
		return t.accessTokenSecret, nil
	})

	if err != nil {
		// JWT v5에서는 에러 타입이 변경됨
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, securityError.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, securityError.ErrTokenInvalid
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, securityError.ErrTokenInvalid
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, securityError.ErrTokenInvalid
		}
		return nil, fmt.Errorf("%w: %v", securityError.ErrTokenInvalid, err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, securityError.ErrTokenClaimsParseFailed
	}

	if !token.Valid {
		return nil, securityError.ErrTokenInvalid
	}

	return claims, nil
}

func (t *TokenManager) GetRefreshTokenTTL() time.Duration {
	return t.refreshTokenTTL
}
