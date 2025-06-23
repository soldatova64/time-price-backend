package requests

type UserRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"` // Здесь важно именно 'email'
	Password string `json:"password" validate:"required,min=6"`
}
