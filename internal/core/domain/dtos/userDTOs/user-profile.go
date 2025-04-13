package userDTOs

import "time"

type UserProfile struct {
	User    UserDTO    `json:"user"`
	Session SessionDTO `json:"session"`
	Meta    MetaDTO    `json:"meta"`
}

type MetaDTO struct {
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}
