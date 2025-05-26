package responses

import (
	"main/models"
)

type HomeResponse struct {
	Meta models.Meta `json:"meta"`
	Data interface{} `json:"data"`
}
