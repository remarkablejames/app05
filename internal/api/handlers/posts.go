package handlers

import (
	"app05/internal/core/application/contracts"
	"app05/internal/infrastructure/services"
	"app05/pkg/appErrors"
	"app05/pkg/utils"
	"net/http"
)

type PostHandler struct {
	postService *services.PostService
	logger      contracts.Logger
}

func NewPostHandler(postService *services.PostService, logger contracts.Logger) *PostHandler {
	return &PostHandler{
		postService: postService,
		logger:      logger,
	}
}

// GetAllPosts  is a function that returns all posts
func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	// Get all posts from the posts service
	posts, err := h.postService.GetAllPosts(r.Context())
	if err != nil {
		appError := appErrors.New(appErrors.CodeInternal, "Failed to retrieve posts")
		appErrors.HandleError(w, appError, h.logger)
		return
	}
	// Return the health response as a JSON response
	utils.SendJSON(w, posts)
}
