package thing

import (
	"database/sql"
	"fmt"
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

func (r *Repository) FindAll() ([]entity.Thing, error) {
	rows, err := r.db.Query("SELECT id, name, pay_date, pay_price, sale_date, sale_price FROM thing WHERE is_deleted = FALSE")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	things := []entity.Thing{}

	for rows.Next() {
		var t entity.Thing

		err = rows.Scan(&t.ID, &t.Name, &t.PayDate, &t.PayPrice, &t.SaleDate, &t.SalePrice)

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

func (r *Repository) Add(db *sql.DB, thing *entity.Thing) (*entity.Thing, error) {

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

	query := `INSERT INTO thing(name, pay_date, pay_price, sale_date, sale_price) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.QueryRow(
		query,
		thing.Name,
		thing.PayDate,
		thing.PayPrice,
		saleDate,
		salePrice,
	).Scan(&thing.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return thing, nil
}

func (r *Repository) Update(db *sql.DB, thing *entity.Thing) (*entity.Thing, error) {
	query := "UPDATE thing SET "
	params := []interface{}{}
	paramCount := 1

	if thing.Name != "" {
		query += fmt.Sprintf("name = $%d, ", paramCount)
		params = append(params, thing.Name)
		paramCount++
	}

	if !thing.PayDate.IsZero() {
		query += fmt.Sprintf("pay_date = $%d, ", paramCount)
		params = append(params, thing.PayDate)
		paramCount++
	}

	if thing.PayPrice > 0 {
		query += fmt.Sprintf("pay_price = $%d, ", paramCount)
		params = append(params, thing.PayPrice)
		paramCount++
	}

	if thing.SaleDate.Valid {
		query += fmt.Sprintf("sale_date = $%d, ", paramCount)
		params = append(params, thing.SaleDate.Time)
		paramCount++
	}

	if thing.SalePrice.Valid {
		query += fmt.Sprintf("sale_price = $%d, ", paramCount)
		params = append(params, thing.SalePrice.Int64)
		paramCount++
	}
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
		thing.ID)

	query = strings.TrimSuffix(query, ", ")

	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, thing.ID)

	_, err = r.db.Exec(query, params...)
	if err != nil {
		return nil, err
	}

	return thing, nil
}
