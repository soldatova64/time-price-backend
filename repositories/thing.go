package repositories

import (
	"database/sql"
	"main/entity"
	"time"
)

func FindAll(db *sql.DB) ([]entity.Thing, error) {
	rows, err := db.Query("SELECT id, name, pay_date, pay_price, sale_date, sale_price FROM thing")

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

func SoftDelete(db *sql.DB, id int) error {
	query := `UPDATE thing SET deleted=true, deleted_at=$1, WHERE id=$2`
	_, err := db.Exec(query, time.Now(), id)
	return err
}

func FindByID(db *sql.DB, id int) (*entity.Thing, error) {
	query := `SELECT id, name, pay_date, pay_price, sale_date, sale_price FROM thing WHERE id=$1 and deleted=false and deleted_at IS NULL`
	var t entity.Thing
	err := db.QueryRow(query, id).Scan(
		&t.ID,
		&t.Name,
		&t.PayDate,
		&t.PayPrice,
		&t.SaleDate,
		&t.SalePrice,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
