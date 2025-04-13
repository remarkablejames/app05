package handlers

import (
	_ "app05/internal/core/application/constants"
	"app05/internal/core/application/contracts"
	"app05/internal/core/domain/dtos/userDTOs"
	"app05/internal/infrastructure/services"
	"app05/pkg/appErrors"
	"app05/pkg/utils"
	"github.com/go-playground/validator/v10"
	_ "github.com/google/uuid"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
	validator   *validator.Validate
	logger      contracts.Logger
}

func NewAuthHandler(authService *services.AuthService, logger contracts.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
		logger:      logger,
	}
}

// Request payload for forgot password
type forgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Request payload for resetting the password
type resetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var input userDTOs.RegisterUserRequest
	if err := utils.ParseJSON(w, r, &input); err != nil {
		appError := appErrors.New(appErrors.CodeBadRequest, err.Error())
		appErrors.HandleError(w, appError, h.logger)
		return
	}

	if err := h.validator.Struct(input); err != nil {
		appError := appErrors.New(appErrors.CodeBadRequest, err.Error())
		appErrors.HandleError(w, appError, h.logger)
		return
	}

	_, err := h.authService.Register(r.Context(), input)

	if err != nil {
		appErrors.HandleError(w, err, h.logger)
		return
	}

	// auto login user after registration
	loginInput := userDTOs.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	}

	var response *userDTOs.LoginResponse
	response, err = h.authService.Login(r.Context(), loginInput)

	if err != nil {
		appErrors.HandleError(w, err, h.logger)
		return
	}

	utils.SendJSON(w, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var input userDTOs.LoginInput
	if err := utils.ParseJSON(w, r, &input); err != nil {
		appErr := appErrors.New(appErrors.CodeBadRequest, err.Error())
		appErrors.HandleError(w, appErr, h.logger)
		return
	}

	input.DeviceInfo.IPAddress = r.Header.Get("X-Real-IP")
	if input.DeviceInfo.IPAddress == "" {
		input.DeviceInfo.IPAddress = r.RemoteAddr
	}

	if err := h.validator.Struct(input); err != nil {
		appErr := appErrors.New(appErrors.CodeBadRequest, err.Error())
		appErrors.HandleError(w, appErr, h.logger)
		return
	}

	response, err := h.authService.Login(r.Context(), input)
	if err != nil {
		appErrors.HandleError(w, err, h.logger)
		return
	}

	utils.SendJSON(w, response)
}
