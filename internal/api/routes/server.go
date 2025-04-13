package routes

import (
	"app05/internal/api/handlers"
	"app05/internal/core/application/contracts"
	"app05/internal/infrastructure/services"
	"github.com/go-chi/chi/v5"
)

func RegisterServerStatusRoutes(r chi.Router, healthService *services.ServerService, logger contracts.Logger) {
	h := handlers.NewServerStatusHandler(healthService)
	r.Route("/server", func(r chi.Router) {
		r.Get("/health", h.HealthCheck)
	})
}
