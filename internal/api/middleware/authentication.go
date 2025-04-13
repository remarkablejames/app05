package middlewares

import (
	"app05/internal/core/application/constants"
	"app05/internal/core/application/contracts"
	"app05/internal/infrastructure/cache"
	"app05/pkg/appErrors"
	"context"
	"net/http"
)

func NewAuthMiddleware(sessionCache *cache.SessionCache, logger contracts.Logger) *AuthMiddlewareImpl {
	return &AuthMiddlewareImpl{
		sessionCache: sessionCache,
		logger:       logger,
	}
}

type AuthMiddlewareImpl struct {
	sessionCache *cache.SessionCache
	logger       contracts.Logger
}

// RequireAuth - middleware that requires authentication
func (a *AuthMiddlewareImpl) RequireAuth(next http.Handler) http.Handler {
	return a.HandleAuth(next, false)
}

// OptionalAuth - middleware that makes authentication optional
func (a *AuthMiddlewareImpl) OptionalAuth(next http.Handler) http.Handler {
	return a.HandleAuth(next, true)
}

// HandleAuth - core authentication logic with optional flag
func (a *AuthMiddlewareImpl) HandleAuth(next http.Handler, optional bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)

		// If no token and authentication is optional, proceed without user context
		if token == "" {
			if optional {
				next.ServeHTTP(w, r)
				return
			}

			// Otherwise, return unauthorized error
			Error := appErrors.New(appErrors.CodeUnauthorized, "Authorization header was not provided")
			appErrors.HandleError(w, Error, a.logger)
			return
		}

		// If token exists, validate it regardless of optional flag
		valid, err := a.sessionCache.ValidateSession(r.Context(), token)
		if err != nil || !valid {
			Error := appErrors.New(appErrors.CodeUnauthorized, "Invalid or expired session")
			appErrors.HandleError(w, Error, a.logger)
			return
		}

		// Health check
		if !a.sessionCache.IsSessionHealthy(r.Context(), token) {
			Error := appErrors.New(appErrors.CodeUnauthorized, "Session expired or revoked")
			appErrors.HandleError(w, Error, a.logger)
			return
		}

		session, err := a.sessionCache.GetSession(r.Context(), token)
		if err != nil {
			Error := appErrors.New(appErrors.CodeUnauthorized, "Invalid or expired session")
			appErrors.HandleError(w, Error, a.logger)
			return
		}

		// Add session and user role to context
		ctx := context.WithValue(r.Context(), constants.SessionCtxKey, session)
		ctx = context.WithValue(ctx, constants.UserRoleCtxKey, session.UserRole)
		ctx = context.WithValue(ctx, constants.UserIdCtxKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
