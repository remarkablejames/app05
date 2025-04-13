package handlers

import (
	"app05/internal/core/application/constants"
	"app05/internal/core/application/contracts"
	"app05/internal/infrastructure/cache"
	"app05/internal/infrastructure/services"
	"app05/pkg/appErrors"
	"app05/pkg/utils"
	"github.com/google/uuid"
	"net/http"
)

const (
	maxUploadSize = 5 << 20 // 5MB
	uploadPath    = "uploads/profile-pictures"
)

type UserHandler struct {
	logger      contracts.Logger
	userService *services.UserService
	cache       *cache.SessionCache
}

func NewUserHandler(logger contracts.Logger, userService *services.UserService, cache *cache.SessionCache) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
		cache:       cache,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIdCtxKey).(uuid.UUID)
	if !ok {
		appError := appErrors.New(appErrors.CodeUnauthorized, "login required")
		appErrors.HandleError(w, appError, h.logger)
		return
	}

	profile, err := h.userService.GetUserProfile(ctx, userId)
	if err != nil {
		appErrors.HandleError(w, err, h.logger)
		return
	}

	utils.SendJSON(w, profile)
}

func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, ok := ctx.Value(constants.UserIdCtxKey).(uuid.UUID)
	if !ok {
		appError := appErrors.New(appErrors.CodeUnauthorized, "authentication required")
		appErrors.HandleError(w, appError, h.logger)
		return
	}

	// Parse request body
	var input struct {
		VerificationCode string `json:"code" validate:"required"`
	}
	if err := utils.ParseJSON(w, r, &input); err != nil {
		appErrors.HandleError(w, err, h.logger)
		return
	}

	// Update user profile
	err := h.userService.VerifyEmail(ctx, userID, input.VerificationCode)
	if err != nil {
		appErrors.HandleError(w, err, h.logger)
		return
	}

	// Return success message
	utils.SendJSON(w, map[string]string{"message": "User email has been verified"})
}
