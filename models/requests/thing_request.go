package requests

import (
	"main/types"
	"time"
)

type ThingRequest struct {
	Name      string          `json:"name" validate:"required,min=3"`
	PayDate   time.Time       `json:"pay_date" validate:"required"`
	PayPrice  int             `json:"pay_price" validate:"required,gt=0"`
	SaleDate  types.NullTime  `json:"sale_date"`
	SalePrice types.NullInt64 `json:"sale_price"`
	UserID    int             `json:"user_id"`
}
