package entity

import (
	"database/sql"
	"time"
)

type Thing struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	PayDate   time.Time     `json:"pay_date"`
	PayPrice  int           `json:"pay_price"`
	SaleDate  sql.NullTime  `json:"sale_date"`
	SalePrice sql.NullInt64 `json:"sale_price"`
	Days      int           `json:"days"`
	PayDay    float64       `json:"pay_day"`
	CreatedAt string        `json:"-"`
	Deleted   bool          `json:"-"`
	DeletedAt sql.NullTime  `json:"-"`
}
