package auth_token

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"main/entity"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddToken(authToken *entity.AuthToken, durationHours int) (*entity.AuthToken, error) {
	endDate := time.Now().Add(time.Duration(durationHours) * time.Hour)

	query := `INSERT INTO auth_tokens (user_id, token, end_date) 
              VALUES ($1, $2, $3) RETURNING id, created_at`

	err := r.db.QueryRow(
		query,
		authToken.UserID,
		authToken.Token,
		endDate,
	).Scan(&authToken.ID,
		&authToken.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return authToken, nil
}

func (r *Repository) FindByToken(token string) (entity.AuthToken, error) {
	if token == "" {
		return entity.AuthToken{}, fmt.Errorf("пустой токен")
	}

	query := `SELECT id, user_id, token, created_at, end_date 
              FROM auth_tokens WHERE token = $1`

	var authToken entity.AuthToken
	err := r.db.QueryRow(query, token).Scan(
		&authToken.ID,
		&authToken.UserID,
		&authToken.Token,
		&authToken.CreatedAt,
		&authToken.EndDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.AuthToken{}, nil
		}
		log.Printf("Ошибка при поиске токена: %v\nЗапрос: %s\nТокен: %s", err, query, token)
		return entity.AuthToken{}, fmt.Errorf("ошибка базы данных: %w", err)
	}

	return authToken, nil
}
