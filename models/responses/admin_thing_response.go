package responses

import (
	"main/entity"
	"main/models"
)

type AdminThingResponse struct {
	Meta models.Meta  `json:"meta"`
	Data entity.Thing `json:"data,omitempty"`
}
