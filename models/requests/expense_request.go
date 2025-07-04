package requests

import "time"

type ExpenseRequest struct {
	ThingID     int       `json:"thing_id" validate:"required,gt=0"`
	Sum         int       `json:"sum" validate:"required,gt=0"`
	Description string    `json:"description" validate:"required,min=3"`
	ExpenseDate time.Time `json:"expense_date" validate:"required"`
}
