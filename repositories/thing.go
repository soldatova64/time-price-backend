package repositories

import (
	"database/sql"
	"main/entity"
)

func FindAll(db *sql.DB) ([]entity.Thing, error) {
	rows, err := db.Query("SELECT id, name, pay_date, pay_price, sale_date, sale_price FROM thing WHERE is_deleted = FALSE")

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
