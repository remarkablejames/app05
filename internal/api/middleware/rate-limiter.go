package middlewares

import (
	"app05/internal/infrastructure/rate_limiter"
	"net/http"
)

func RateLimiterMiddleware(limiter *rate_limiter.FixedWindowRateLimiter, enabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Rate limiter logic
			if enabled {
				if allow, retryAfter := limiter.Allow(r.RemoteAddr); !allow {
					http.Error(w, "Too many requests", http.StatusTooManyRequests)
					w.Header().Set("Retry-After", retryAfter.String())
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
