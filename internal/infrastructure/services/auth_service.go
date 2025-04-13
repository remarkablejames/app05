package services

import (
	"app05/internal/core/application/contracts"
	"app05/internal/core/domain/dtos/userDTOs"
	"app05/internal/core/domain/entities"
	"app05/internal/core/domain/repositories"
	"app05/internal/infrastructure/cache"
	"app05/pkg/appErrors"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	resetTokenLength = 32
	resetTokenExpiry = 1 * time.Hour // Token expires after 1 hour
)

type AuthService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
	sessionMgr  *entities.SessionManager
	cache       *cache.SessionCache
	logger      contracts.Logger
}

func NewAuthService(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	cache *cache.SessionCache,
	logger contracts.Logger,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionMgr:  entities.NewSessionManager(),
		cache:       cache,
		logger:      logger,
	}
}

func (s *AuthService) Register(ctx context.Context, input userDTOs.RegisterUserRequest) (*entities.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, appErrors.New(appErrors.CodeBadRequest, "Email already taken by another user. Try a different email")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Email:                  input.Email,
		HashedPassword:         string(hashedPassword),
		FirstName:              input.FirstName,
		LastName:               input.LastName,
		SubscribedToNewsletter: input.SubscribedToNewsletter,
		Role:                   entities.RoleStudent,
		Active:                 true,
		EmailVerified:          false,
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input userDTOs.LoginInput) (*userDTOs.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		fmt.Println(err)
		return nil, appErrors.New(appErrors.CodeBadRequest, "invalid email or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password))
	if err != nil {
		// TODO: when password is incorrect, the user gets a weird error message
		return nil, appErrors.New(appErrors.CodeBadRequest, "invalid email or password")
	}

	if !user.Active {
		return nil, appErrors.New(appErrors.CodeBadRequest, "Account is Deactivated")
	}

	// Check for existing active session in Redis
	existingToken, err := s.cache.GetUserActiveSession(ctx, user.ID)
	if err == nil && existingToken != "" {
		// Delete existing session in Redis
		err := s.cache.DeleteSession(ctx, existingToken, user.ID)
		if err != nil {
			return nil, appErrors.New(appErrors.CodeInternal, "failed to delete existing session")
		}
		// Update DB asynchronously with error logging
		go func() {
			ctx := context.Background()
			existingSession, err := s.sessionRepo.GetSessionByToken(ctx, existingToken)
			if err != nil || existingSession == nil {
				// Optionally log or ignore if session is absent
				return
			}
			existingSession.Status = entities.SessionStatusRevoked
			if err := s.sessionRepo.UpdateSession(ctx, existingSession); err != nil {
				s.logger.Error("Failed to update session status", "error", err)
			}
		}()
	}

	// Generate tokens
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	// Create new session
	session := &entities.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		Token:        token,
		UserRole:     user.Role,
		RefreshToken: refreshToken,
		Status:       entities.SessionStatusActive,
		DeviceInfo:   input.DeviceInfo,
		ExpiresAt:    time.Now().Add(s.sessionMgr.SessionDuration),
	}

	// Store session in Redis
	if err := s.cache.StoreSession(ctx, session); err != nil {
		return nil, err
	}

	// Store in DB asynchronously with error logging
	go func() {
		ctx := context.Background()
		if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
			s.logger.Error("Failed to create session in database", "error", err)
		}
	}()

	user.SetCurrentSession(session)

	// Create response DTO
	response := &userDTOs.LoginResponse{
		User: userDTOs.UserDTO{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			ProfilePictureURL: func() string {
				if user.ProfilePictureURL != nil {
					return *user.ProfilePictureURL
				}
				return ""
			}(),
			Role:          user.Role.String(),
			Active:        user.Active,
			EmailVerified: user.EmailVerified,
		},
		Session: userDTOs.SessionDTO{
			Token:        session.Token,
			RefreshToken: session.RefreshToken,
			ExpiresAt:    session.ExpiresAt.Format(time.RFC3339),
		},
	}

	return response, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
