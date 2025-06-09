package responses

import (
	"main/entity"
	"main/models"
)

type AdminExpenseResponse struct {
	Meta models.Meta    `json:"meta"`
	Data entity.Expense `json:"data"`
}
