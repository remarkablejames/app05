package userDTOs

type UserDTO struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Role              string `json:"role"`
	Active            bool   `json:"active"`
	EmailVerified     bool   `json:"email_verified"`
}

type SessionDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}
