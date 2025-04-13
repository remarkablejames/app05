package userDTOs

type LoginResponse struct {
	User    UserDTO    `json:"user"`
	Session SessionDTO `json:"session"`
}
