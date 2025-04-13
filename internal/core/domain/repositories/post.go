package repositories

import (
	"app05/internal/core/domain/dtos/postDTOs"
	"context"
)

type PostRepository interface {
	// GetAllPosts retrieves all posts
	GetAllPosts(ctx context.Context) ([]*postDTOs.PostDTO, error)
}
