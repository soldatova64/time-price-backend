package requests

type UserUpdateRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
}
