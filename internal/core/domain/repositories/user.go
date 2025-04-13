package repositories

import (
	"app05/internal/core/domain/entities"
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetUserProfile(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetUserByResetToken(ctx context.Context, hashedToken string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	UpdateProfilePicture(ctx context.Context, userID uuid.UUID, removeProfilePicture bool, profilePictureURL string) error

	CreateVerificationCode(ctx context.Context, userID uuid.UUID, code string, expiresAt time.Time) error
	VerifyEmail(ctx context.Context, userID uuid.UUID, code string) error
}
