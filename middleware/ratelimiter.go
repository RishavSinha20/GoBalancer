package middleware

import (
	"net/http"
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(limit int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, limit),
	}

	for i := 0; i < limit; i++ {
		rl.tokens <- struct{}{}
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return rl
}
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-rl.tokens:
			next.ServeHTTP(w, r)

		default:
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	})
}
