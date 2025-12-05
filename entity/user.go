package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int          `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Password  string       `json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	IsDeleted bool         `json:"is_deleted"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
