package entity

import (
	"database/sql"
)

type Expense struct {
	ID          int          `json:"id"`
	ThingID     int          `json:"thing_id"`
	Sum         int          `json:"sum"`
	Description string       `json:"description"`
	CreatedAt   string       `json:"-"`
	Deleted     bool         `json:"-"`
	DeletedAt   sql.NullTime `json:"-"`
}
