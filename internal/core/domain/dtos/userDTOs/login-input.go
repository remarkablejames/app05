package userDTOs

import "app05/internal/core/domain/entities"

type LoginInput struct {
	Email      string              `json:"email" validate:"required,email"`
	Password   string              `json:"password" validate:"required"`
	DeviceInfo entities.DeviceInfo `json:"device_info" validate:"required"`
}
