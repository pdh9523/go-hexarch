package jwt

import (
	"errors"
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
		return nil, securityError.ErrTokenGenerationFailed
	}
	refreshToken, err := t.generateRefreshToken(userID, role)
	if err != nil {
		return nil, securityError.ErrTokenGenerationFailed
	}

	return &Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (t *TokenManager) ValidateAccessToken(accessToken string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			if _, ok := token.Header["alg"].(string); !ok {
				return nil, securityError.ErrTokenMissingAlgorithm
			}
			return nil, securityError.ErrTokenUnexpectedSigningMethod
		}
		return t.accessTokenSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, securityError.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) || errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, securityError.ErrTokenInvalid
		}
		return nil, securityError.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, securityError.ErrTokenClaimsParseFailed
	}

	if claims.TokenType != "access_token" {
		return nil, securityError.ErrTokenInvalid
	}
	return claims, nil
}

func (t *TokenManager) ValidateRefreshToken(refreshToken string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			if _, ok := token.Header["alg"].(string); !ok {
				return nil, securityError.ErrTokenMissingAlgorithm
			}
			return nil, securityError.ErrTokenUnexpectedSigningMethod
		}
		return t.refreshTokenSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, securityError.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) || errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, securityError.ErrTokenInvalid
		}
		return nil, securityError.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, securityError.ErrTokenClaimsParseFailed
	}
	if claims.TokenType != "refresh_token" {
		return nil, securityError.ErrTokenInvalid
	}
	return claims, nil
}

func (t *TokenManager) RefreshTokens(refreshToken string) (*Token, error) {
	claims, err := t.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
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
		return "", securityError.ErrTokenInvalid
	}
	return bearerToken[7:], nil
}

func (t *TokenManager) ExtractAccessTokenClaims(accessToken string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 알고리즘 검증 - 화이트리스트 방식
		if token.Method != jwt.SigningMethodHS256 {
			if _, ok := token.Header["alg"].(string); !ok {
				return nil, securityError.ErrTokenMissingAlgorithm
			}
			return nil, securityError.ErrTokenUnexpectedSigningMethod
		}
		return t.accessTokenSecret, nil
	})

	if err != nil {
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
		return nil, securityError.ErrTokenInvalid
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
