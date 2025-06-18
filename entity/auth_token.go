package entity

import "time"

type AuthToken struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	EndDate   time.Time `json:"end_date"`
}
