package routes

import (
	"app05/internal/api/handlers"
	"app05/internal/core/application/contracts"
	"app05/internal/infrastructure/services"
	"github.com/go-chi/chi/v5"
)

func RegisterPostRoutes(r chi.Router, postService *services.PostService, logger contracts.Logger) {
	h := handlers.NewPostHandler(postService, logger)
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", h.GetAllPosts)
	})
}
