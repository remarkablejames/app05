package userDTOs

type RegisterUserRequest struct {
	Email                  string `json:"email" validate:"required,email"`
	Password               string `json:"password" validate:"required,min=8"`
	FirstName              string `json:"first_name" validate:"required"`
	LastName               string `json:"last_name" validate:"required"`
	SubscribedToNewsletter bool   `json:"subscribed_to_newsletter" validate:"required"`
}
