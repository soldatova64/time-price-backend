package responses

import (
	"main/entity"
	"main/models"
)

type HomeResponse struct {
	Meta models.Meta    `json:"meta"`
	Data []entity.Thing `json:"data"`
}
