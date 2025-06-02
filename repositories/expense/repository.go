package expense

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
func (r *Repository) FindAll() ([]entity.Expense, error) {
	rows, err := r.db.Query("SELECT id, thing_id, sum, description, expense_date, created_at  FROM expense WHERE is_deleted = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := []entity.Expense{}

	for rows.Next() {
		var e entity.Expense
		err = rows.Scan(&e.ID, &e.ThingID, &e.Sum, &e.Description, &e.ExpenseDate, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}
