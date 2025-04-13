package services

import (
	"app05/internal/core/application/contracts"
	"app05/internal/core/domain/entities"
	"app05/internal/core/domain/repositories"
	"context"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo repositories.UserRepository
	logger   contracts.Logger
}

func NewUserService(userRepo repositories.UserRepository, logger contracts.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserProfile(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetUserProfile(ctx, id)
}

func (s *UserService) VerifyEmail(ctx context.Context, userID uuid.UUID, code string) error {
	return s.userRepo.VerifyEmail(ctx, userID, code)
}
