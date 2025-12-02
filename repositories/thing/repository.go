package thing

import (
	"database/sql"
	"log"
	"main/entity"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(userID int) ([]entity.Thing, error) {
	rows, err := r.db.Query(
		"SELECT id, name, pay_date, pay_price, sale_date, sale_price, user_id FROM thing WHERE is_deleted = FALSE AND user_id = $1",
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	things := []entity.Thing{}

	for rows.Next() {
		var t entity.Thing

		err = rows.Scan(&t.ID, &t.Name, &t.PayDate, &t.PayPrice, &t.SaleDate, &t.SalePrice, &t.UserID)

		if err != nil {
			return nil, err
		}

		things = append(things, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return things, nil
}

func (r *Repository) Add(thing *entity.Thing) (*entity.Thing, error) {
	thing.ID = 0
	payDate := thing.PayDate.Format("2006-01-02")
	var saleDate interface{}
	if thing.SaleDate.Valid {
		saleDate = thing.SaleDate.Time.Format("2006-01-02") // Форматируем дату продажи
	} else {
		saleDate = nil
	}

	var salePrice interface{}
	if thing.SalePrice.Valid {
		salePrice = thing.SalePrice.Int64
	} else {
		salePrice = nil
	}

	query := `INSERT INTO thing(name, pay_date, pay_price, sale_date, sale_price, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(
		query,
		thing.Name,
		payDate,
		thing.PayPrice,
		saleDate,
		salePrice,
		thing.UserID,
	).Scan(&thing.ID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Println("Duplicate key error detected, resetting sequence...")
			_, resetErr := r.db.Exec("SELECT setval('thing_id_seq', (SELECT COALESCE(MAX(id), 1) FROM thing))")
			if resetErr != nil {
				log.Printf("Error resetting sequence: %v", resetErr)
				return nil, err
			}
			return r.Add(thing)
		}
		log.Printf("Database error in Thing Add: %v", err)
		return nil, err
	}

	log.Printf("Successfully created thing with ID: %d", thing.ID)
	return thing, nil
}

func (r *Repository) Find(id, userID int) (entity.Thing, error) {
	rows, err := r.db.Query(
		"SELECT id, name, pay_date, pay_price, sale_date, sale_price, user_id FROM thing WHERE id = $1 AND user_id = $2 AND is_deleted = FALSE",
		id, userID,
	)

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	thing := entity.Thing{}

	for rows.Next() {
		var t entity.Thing

		err = rows.Scan(&t.ID, &t.Name, &t.PayDate, &t.PayPrice, &t.SaleDate, &t.SalePrice, &t.UserID)

		if err != nil {
			log.Println(err)
		}

		thing = t
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
	}
	return thing, nil

}

func (r *Repository) Update(thing entity.Thing) (entity.Thing, error) {
	query := `UPDATE thing SET 
		name = $1, 
		pay_date = $2, 
		pay_price = $3, 
		sale_date = $4, 
		sale_price = $5 
	WHERE id = $6 AND user_id = $7 AND is_deleted = FALSE`

	var saleDate interface{}
	if thing.SaleDate.Valid {
		saleDate = thing.SaleDate.Time
	} else {
		saleDate = nil
	}

	var salePrice interface{}
	if thing.SalePrice.Valid {
		salePrice = thing.SalePrice.Int64
	} else {
		salePrice = nil
	}

	_, err := r.db.Exec(query,
		thing.Name,
		thing.PayDate,
		thing.PayPrice,
		saleDate,
		salePrice,
		thing.ID,
		thing.UserID)
	if err != nil {
		log.Println(err)
	}

	return thing, nil
}

func (r *Repository) Delete(id, userID int) error {
	query := `UPDATE thing SET is_deleted = TRUE, deleted_at = NOW() WHERE id = $1 AND user_id = $2 AND is_deleted = FALSE`

	result, err := r.db.Exec(query, id, userID)
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
