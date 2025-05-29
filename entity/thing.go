package entity

import (
	"database/sql"
	"main/types"
	"time"
)

type Thing struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	PayDate   time.Time       `json:"pay_date"`
	PayPrice  int             `json:"pay_price"`
	SaleDate  types.NullTime  `json:"sale_date"`
	SalePrice types.NullInt64 `json:"sale_price"`
	Days      int             `json:"days"`
	PayDay    float64         `json:"pay_day"`
	CreatedAt string          `json:"-"`
	Deleted   bool            `json:"-"`
	DeletedAt sql.NullTime    `json:"-"`
	Expense   []Expense       `json:"expense"`
}
