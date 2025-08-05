package middleware

import (
	"errors"
	securityError "github.com/pdh9523/go-hexarch/shared/security/error_code"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	"github.com/pdh9523/go-hexarch/shared/security/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Security Manager
type MockSecurityManager struct {
	mock.Mock
}

func (m *MockSecurityManager) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockSecurityManager) VerifyPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockSecurityManager) GenerateTokens(username, role string) (*jwt.Token, error) {
	args := m.Called(username, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *MockSecurityManager) ValidateAccessToken(accessToken string) (*jwt.AccessTokenClaims, error) {
	args := m.Called(accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.AccessTokenClaims), args.Error(1)
}

func (m *MockSecurityManager) ValidateRefreshToken(refreshToken string) (*jwt.RefreshTokenClaims, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.RefreshTokenClaims), args.Error(1)
}

func (m *MockSecurityManager) RefreshTokens(refreshToken string) (*jwt.Token, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *MockSecurityManager) ExtractTokenFromBearer(bearerToken string) (string, error) {
	args := m.Called(bearerToken)
	return args.String(0), args.Error(1)
}

func (m *MockSecurityManager) GetRefreshTokenTTL() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockSecurityManager) ExtractAccessTokenClaims(accessToken string) (*jwt.AccessTokenClaims, error) {
	args := m.Called(accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.AccessTokenClaims), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given no authorization header when processing request then should abort with error", func(t *testing.T) {
		responder := response.NewResponder()
		securityManager := new(MockSecurityManager)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := AuthMiddleware(responder, securityManager)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Given invalid bearer format when processing request then should abort with error", func(t *testing.T) {
		responder := response.NewResponder()
		securityManager := new(MockSecurityManager)

		securityManager.On("ExtractTokenFromBearer", "InvalidBearerFormat").Return("", securityError.ErrTokenInvalid)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "InvalidBearerFormat")

		middleware := AuthMiddleware(responder, securityManager)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		securityManager.AssertExpectations(t)
	})

	t.Run("Given empty token when processing request then should abort with error", func(t *testing.T) {
		responder := response.NewResponder()
		securityManager := new(MockSecurityManager)

		securityManager.On("ExtractTokenFromBearer", "Bearer ").Return("", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer ")

		middleware := AuthMiddleware(responder, securityManager)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		securityManager.AssertExpectations(t)
	})

	t.Run("Given invalid token when processing request then should abort with error", func(t *testing.T) {
		responder := response.NewResponder()
		securityManager := new(MockSecurityManager)

		securityManager.On("ExtractTokenFromBearer", "Bearer invalid_token").Return("invalid_token", nil)
		var nilClaims *jwt.AccessTokenClaims
		securityManager.On("ExtractAccessTokenClaims", "invalid_token").Return(nilClaims, errors.New("invalid token"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid_token")

		middleware := AuthMiddleware(responder, securityManager)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		securityManager.AssertExpectations(t)
	})

	t.Run("Given valid token when processing request then should set user context and continue", func(t *testing.T) {
		responder := response.NewResponder()
		securityManager := new(MockSecurityManager)

		claims := &jwt.AccessTokenClaims{
			StandardClaims: jwt.StandardClaims{
				UserID: "user123",
				Role:   "admin",
			},
			TokenType: "access_token",
		}

		securityManager.On("ExtractTokenFromBearer", "Bearer valid_token").Return("valid_token", nil)
		securityManager.On("ExtractAccessTokenClaims", "valid_token").Return(claims, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer valid_token")

		middleware := AuthMiddleware(responder, securityManager)
		middleware(c)

		assert.False(t, c.IsAborted())

		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, "user123", userID)

		role, exists := c.Get("role")
		assert.True(t, exists)
		assert.Equal(t, "admin", role)

		securityManager.AssertExpectations(t)
	})
}

func TestOptionalAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given no authorization header when processing request then should continue without setting context", func(t *testing.T) {
		securityManager := new(MockSecurityManager)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := OptionalAuthMiddleware(securityManager)
		middleware(c)

		assert.False(t, c.IsAborted())

		_, userExists := c.Get("user_id")
		assert.False(t, userExists)

		_, roleExists := c.Get("role")
		assert.False(t, roleExists)
	})

	t.Run("Given invalid token when processing request then should continue without setting context", func(t *testing.T) {
		securityManager := new(MockSecurityManager)

		securityManager.On("ExtractTokenFromBearer", "Bearer invalid_token").Return("invalid_token", nil)
		var nilClaims *jwt.AccessTokenClaims
		securityManager.On("ExtractAccessTokenClaims", "invalid_token").Return(nilClaims, errors.New("invalid token"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid_token")

		middleware := OptionalAuthMiddleware(securityManager)
		middleware(c)

		assert.False(t, c.IsAborted())

		_, userExists := c.Get("user_id")
		assert.False(t, userExists)

		_, roleExists := c.Get("role")
		assert.False(t, roleExists)

		securityManager.AssertExpectations(t)
	})

	t.Run("Given valid token when processing request then should set user context and continue", func(t *testing.T) {
		securityManager := new(MockSecurityManager)

		claims := &jwt.AccessTokenClaims{
			StandardClaims: jwt.StandardClaims{
				UserID: "user123",
				Role:   "user",
			},
			TokenType: "access_token",
		}

		securityManager.On("ExtractTokenFromBearer", "Bearer valid_token").Return("valid_token", nil)
		securityManager.On("ExtractAccessTokenClaims", "valid_token").Return(claims, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer valid_token")

		middleware := OptionalAuthMiddleware(securityManager)
		middleware(c)

		assert.False(t, c.IsAborted())

		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, "user123", userID)

		role, exists := c.Get("role")
		assert.True(t, exists)
		assert.Equal(t, "user", role)

		securityManager.AssertExpectations(t)
	})

	t.Run("Given bearer extraction error when processing request then should continue without setting context", func(t *testing.T) {
		securityManager := new(MockSecurityManager)

		securityManager.On("ExtractTokenFromBearer", "InvalidFormat").Return("", errors.New("invalid format"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "InvalidFormat")

		middleware := OptionalAuthMiddleware(securityManager)
		middleware(c)

		assert.False(t, c.IsAborted())

		_, userExists := c.Get("user_id")
		assert.False(t, userExists)

		_, roleExists := c.Get("role")
		assert.False(t, roleExists)

		securityManager.AssertExpectations(t)
	})
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given request when processing then should set all security headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := SecurityHeadersMiddleware()
		middleware(c)

		assert.False(t, c.IsAborted())
		assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
		assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
	})
}
