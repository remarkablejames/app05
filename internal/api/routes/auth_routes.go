package routes

import (
	"app05/internal/api/handlers"
	middlewares "app05/internal/api/middleware"
	"app05/internal/core/application/contracts"
	"app05/internal/infrastructure/cache"
	"app05/internal/infrastructure/services"
	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, sessionCache *cache.SessionCache, authService *services.AuthService, logger contracts.Logger) {
	h := handlers.NewAuthHandler(authService, logger)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)

		// Protected routes group
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(sessionCache, logger))
		})
	})
}
