package routes

import (
	"app05/internal/api/handlers"
	middlewares "app05/internal/api/middleware"
	"app05/internal/core/application/contracts"
	"app05/internal/infrastructure/cache"
	"app05/internal/infrastructure/services"
	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, sessionCache *cache.SessionCache, userService *services.UserService, logger contracts.Logger) {
	// Create handlers
	userHandler := handlers.NewUserHandler(logger, userService, sessionCache)

	// Protected routes group
	r.Group(func(r chi.Router) {
		// Apply auth middleware
		r.Use(middlewares.AuthMiddleware(sessionCache, logger))

		// User routes
		r.Route("/users", func(r chi.Router) {
			// used to list all users in the CMS. Only accessible by admins and superusers. superusers can list all users, admins can only list users with roles below them
			r.Get("/profile", userHandler.GetProfile)
			r.Post("/verify-email", userHandler.VerifyEmail)
		})
	})
}
