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

func (r *Repository) Add(expense *entity.Expense) (*entity.Expense, error) {
	query := `INSERT INTO expense(thing_id, sum, description, expense_date) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id`

	err := r.db.QueryRow(
		query,
		expense.ThingID,
		expense.Sum,
		expense.Description,
		expense.ExpenseDate,
	).Scan(&expense.ID)

	if err != nil {
		return nil, err
	}

	return expense, nil

}

func (r *Repository) Update(expense *entity.Expense) (*entity.Expense, error) {
	query := `UPDATE expense 
	          SET thing_id = $1, sum = $2, description = $3, expense_date = $4 
	          WHERE id = $5 AND is_deleted = FALSE 
	          RETURNING id, thing_id, sum, description, expense_date`

	err := r.db.QueryRow(
		query,
		expense.ThingID,
		expense.Sum,
		expense.Description,
		expense.ExpenseDate,
		expense.ID,
	).Scan(&expense.ID, &expense.ThingID, &expense.Sum, &expense.Description, &expense.ExpenseDate)

	if err != nil {
		return nil, err
	}

	return expense, nil
}

func (r *Repository) Delete(id int) error {
	query := `UPDATE expense SET is_deleted = TRUE, deleted_at = NOW() WHERE id = $1 AND is_deleted = FALSE`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
