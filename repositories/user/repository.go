package user

import (
	"database/sql"
	"errors"
	"main/entity"
	"main/helpers"
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
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *Repository) FindByUsernameAndPassword(username, password string) (*entity.User, error) {
	query := `SELECT id, username, password FROM users 
              WHERE username = $1 AND is_deleted = FALSE`
	row := r.db.QueryRow(query, username)

	var user entity.User
	var hashedPassword string
	err := row.Scan(&user.ID, &user.Username, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if !helpers.CheckPasswordHash(password, hashedPassword) {
		return nil, nil
	}

	return &user, nil
}

func (r *Repository) Add(user *entity.User) (*entity.User, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return nil, errors.New("all fields are required")
	}

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRow(query, user.Username, user.Email, user.Password)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}
