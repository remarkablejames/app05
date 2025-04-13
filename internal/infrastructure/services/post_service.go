package services

import (
	"app05/internal/core/domain/dtos/postDTOs"
	"app05/internal/core/domain/repositories"
	"context"
)

type PostService struct {
	postRepo repositories.PostRepository
}

func NewPostService(postRepo repositories.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

// GetAllPosts retrieves all posts
func (s *PostService) GetAllPosts(ctx context.Context) ([]*postDTOs.PostDTO, error) {
	posts, err := s.postRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
