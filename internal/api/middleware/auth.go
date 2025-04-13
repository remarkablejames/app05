package middlewares

import (
	"app05/internal/core/application/constants"
	"app05/internal/core/application/contracts"
	"app05/internal/core/domain/entities"
	"app05/internal/infrastructure/cache"
	"app05/pkg/appErrors"
	"context"
	"net/http"
	"strings"
)

// RoleMiddleware checks if the user has the required role
func RoleMiddleware(requiredRoles []entities.Role, logger contracts.Logger, cache *cache.SessionCache) func(next http.Handler) http.Handler {
	// Create a map for O(1) role lookup
	allowedRoles := make(map[entities.Role]bool)
	for _, role := range requiredRoles {
		allowedRoles[role] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value(constants.SessionCtxKey).(*entities.Session)
			if session == nil {
				Error := appErrors.New(appErrors.CodeUnauthorized, "Invalid or expired session. Please login again")
				appErrors.HandleError(w, Error, logger)
				return
			}

			// O(1) role check instead of loop
			if !allowedRoles[session.UserRole] {
				Error := appErrors.New(appErrors.CodeForbidden, "You don't have permission to access this resource")
				appErrors.HandleError(w, Error, logger)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware(sessionCache *cache.SessionCache, logger contracts.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				Error := appErrors.New(appErrors.CodeUnauthorized, "Authorization header was not provided")
				appErrors.HandleError(w, Error, logger)
				return
			}

			// Quick validation check
			valid, err := sessionCache.ValidateSession(r.Context(), token)
			if err != nil || !valid {
				Error := appErrors.New(appErrors.CodeUnauthorized, "Invalid or expired session")
				appErrors.HandleError(w, Error, logger)
				return
			}

			// Health check
			if !sessionCache.IsSessionHealthy(r.Context(), token) {
				Error := appErrors.New(appErrors.CodeUnauthorized, "Session expired or revoked")
				appErrors.HandleError(w, Error, logger)
				return
			}

			session, err := sessionCache.GetSession(r.Context(), token)
			if err != nil {
				Error := appErrors.New(appErrors.CodeUnauthorized, "Invalid or expired session")
				appErrors.HandleError(w, Error, logger)
				return
			}

			// Add session and user role to context
			ctx := context.WithValue(r.Context(), constants.SessionCtxKey, session)
			ctx = context.WithValue(ctx, constants.UserRoleCtxKey, session.UserRole)
			ctx = context.WithValue(ctx, constants.UserIdCtxKey, session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
