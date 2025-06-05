package user

import (
	"database/sql"
	"main/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll() ([]entity.User, error) {
	rows, err := r.db.Query("SELECT id, username,  email, password, created_at, updated_at, is_deleted, deleted_at FROM users WHERE is_deleted = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []entity.User{}

	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.IsDeleted, &u.DeletedAt)
		if err != nil {
			return nil, err
			users = append(users, u)
		}
	}
	return users, nil
}
