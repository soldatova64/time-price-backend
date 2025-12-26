package requests

type UserUpdateRequest struct {
	Password *string `json:"password,omitempty" validate:"omitempty,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
}
