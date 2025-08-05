package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given first request when middleware processes then should call next", func(t *testing.T) {
		responder := response.NewResponder()
		rateLimiter := NewRateLimiter(time.Minute, 10)
		defer rateLimiter.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := RateLimitMiddleware(responder, rateLimiter)
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("Given requests within burst limit when middleware processes then should call next", func(t *testing.T) {
		responder := response.NewResponder()
		burstLimit := 5
		rateLimiter := NewRateLimiter(time.Minute, burstLimit)
		defer rateLimiter.Close()

		for i := 0; i < burstLimit; i++ {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			middleware := RateLimitMiddleware(responder, rateLimiter)
			middleware(c)

			assert.False(t, c.IsAborted())
		}
	})

	t.Run("Given requests exceeding burst limit when middleware processes then should abort with error", func(t *testing.T) {
		responder := response.NewResponder()
		burstLimit := 5
		rateLimiter := NewRateLimiter(time.Minute, burstLimit)
		defer rateLimiter.Close()

		// Fill up the burst limit
		for i := 0; i < burstLimit; i++ {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			middleware := RateLimitMiddleware(responder, rateLimiter)
			middleware(c)
		}

		// This request should be rate limited
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := RateLimitMiddleware(responder, rateLimiter)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("Given requests after rate period when middleware processes then should reset counter", func(t *testing.T) {
		responder := response.NewResponder()
		ratePeriod := 100 * time.Millisecond
		burstLimit := 2
		rateLimiter := NewRateLimiter(ratePeriod, burstLimit)
		defer rateLimiter.Close()

		// Fill up the burst limit
		for i := 0; i < burstLimit; i++ {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			middleware := RateLimitMiddleware(responder, rateLimiter)
			middleware(c)
			assert.False(t, c.IsAborted())
		}

		// Wait for rate period to pass
		time.Sleep(ratePeriod + 10*time.Millisecond)

		// This request should be allowed after rate period
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := RateLimitMiddleware(responder, rateLimiter)
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("Given different IPs when middleware processes then should track separately", func(t *testing.T) {
		responder := response.NewResponder()
		burstLimit := 1
		rateLimiter := NewRateLimiter(time.Minute, burstLimit)
		defer rateLimiter.Close()

		// First IP - should be allowed
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		c1.Request = httptest.NewRequest("GET", "/test", nil)
		c1.Request.RemoteAddr = "192.168.1.1:8080"

		middleware := RateLimitMiddleware(responder, rateLimiter)
		middleware(c1)

		assert.False(t, c1.IsAborted())

		// Second IP - should also be allowed even though first IP used up its limit
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Request = httptest.NewRequest("GET", "/test", nil)
		c2.Request.RemoteAddr = "192.168.1.2:8080"

		middleware(c2)

		assert.False(t, c2.IsAborted())
	})
}

func TestIPBasedRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Given request when IP-based middleware processes then should call next", func(t *testing.T) {
		responder := response.NewResponder()
		requestsPerMinute := 10

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := IPBasedRateLimitMiddleware(responder, requestsPerMinute)
		middleware(c)

		assert.False(t, c.IsAborted())
	})
}

func TestRateLimiterConcurrency(t *testing.T) {
	rateLimiter := NewRateLimiter(time.Second, 10)
	defer rateLimiter.Close()

	t.Run("Given concurrent requests when rate limiter processes then should handle safely", func(t *testing.T) {
		concurrency := 50
		var wg sync.WaitGroup
		results := make([]bool, concurrency)

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(index int) {
				defer wg.Done()
				results[index] = rateLimiter.Allow("192.168.1.1")
			}(i)
		}

		wg.Wait()

		// Count allowed and denied requests
		allowed := 0
		for _, result := range results {
			if result {
				allowed++
			}
		}

		// Only burst limit should be allowed
		assert.Equal(t, 10, allowed)
		assert.Equal(t, concurrency-10, concurrency-allowed)
	})
}

func TestRateLimiterCleanup(t *testing.T) {
	t.Run("Given expired visitors when cleanup runs then should remove old entries", func(t *testing.T) {
		// Create rate limiter with short cleanup and expire times for testing
		rateLimiter := &RateLimiter{
			visitors:    make(map[string]*Visitor),
			rate:        time.Minute,
			burst:       10,
			cleanupRate: 10 * time.Millisecond,
			expireAfter: 50 * time.Millisecond,
		}

		// Add a visitor
		rateLimiter.Allow("192.168.1.1")

		// Verify visitor exists
		rateLimiter.mu.RLock()
		initialCount := len(rateLimiter.visitors)
		rateLimiter.mu.RUnlock()
		assert.Equal(t, 1, initialCount)

		// Start cleanup goroutine
		ctx, cancel := context.WithCancel(context.Background())
		rateLimiter.ctx = ctx
		rateLimiter.cancelFunc = cancel
		go rateLimiter.cleanupVisitors()

		// Wait for cleanup to remove expired visitor
		time.Sleep(100 * time.Millisecond)

		// Verify visitor was cleaned up
		rateLimiter.mu.RLock()
		finalCount := len(rateLimiter.visitors)
		rateLimiter.mu.RUnlock()
		assert.Equal(t, 0, finalCount)

		rateLimiter.Close()
	})
}
