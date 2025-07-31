package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	commonError "github.com/pdh9523/go-hexarch/shared/common/error_code"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors    map[string]*Visitor
	mu          sync.RWMutex
	rate        time.Duration
	burst       int
	cleanupRate time.Duration
	expireAfter time.Duration
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

type Visitor struct {
	lastSeen time.Time
	count    int
}

func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())

	rl := &RateLimiter{
		visitors:    make(map[string]*Visitor),
		rate:        rate,
		burst:       burst,
		cleanupRate: time.Minute,
		expireAfter: time.Hour,
		ctx:         ctx,
		cancelFunc:  cancel,
	}

	go rl.cleanupVisitors()

	return rl
}

func (l *RateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	visitor, exists := l.visitors[ip]
	if !exists {
		l.visitors[ip] = &Visitor{
			lastSeen: now,
			count:    1,
		}
		return true
	}

	if now.Sub(visitor.lastSeen) > l.rate {
		visitor.count = 1
		visitor.lastSeen = now
		return true
	}

	if visitor.count < l.burst {
		visitor.count++
		visitor.lastSeen = now
		return true
	}
	return false
}

func (l *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(l.cleanupRate)
	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			l.mu.Lock()
			now := time.Now()
			for ip, visitor := range l.visitors {
				if now.Sub(visitor.lastSeen) > l.expireAfter {
					delete(l.visitors, ip)
				}
			}
			l.mu.Unlock()
		}
	}
}

func (l *RateLimiter) Close() {
	if l.cancelFunc != nil {
		l.cancelFunc()
	}
}

func RateLimitMiddleware(responder *response.Responder, rateLimiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rateLimiter.Allow(ip) {
			responder.AbortWithError(c, commonError.ErrRateLimitExceeded)
			return
		}

		c.Next()
	}
}

func IPBasedRateLimitMiddleware(responder *response.Responder, requestsPerMinute int) gin.HandlerFunc {
	rateLimiter := NewRateLimiter(time.Minute, requestsPerMinute)
	return RateLimitMiddleware(responder, rateLimiter)
}
