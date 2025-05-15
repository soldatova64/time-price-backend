package entity

import (
	"main/types"
)

type Thing struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	PayDate   string           `json:"pay_date"`
	PayPrice  int              `json:"pay_price"`
	SaleDate  types.NullString `json:"sale_date"`
	SalePrice types.NullInt64  `json:"sale_price"`
	CreatedAt string           `json:"-"`
}
