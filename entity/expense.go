package entity

import (
	"database/sql"
	"time"
)

type Expense struct {
	ID          int          `json:"id"`
	ThingID     int          `json:"thing_id"`
	Sum         int          `json:"sum"`
	Description string       `json:"description"`
	ExpenseDate time.Time    `json:"expense_date"`
	CreatedAt   string       `json:"-"`
	Deleted     bool         `json:"-"`
	DeletedAt   sql.NullTime `json:"-"`
}
