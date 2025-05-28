package responses

import "main/models"

type ErrorResponse struct {
	Meta   models.Meta `json:"meta"`
	Errors []Error     `json:"errors"`
}

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
