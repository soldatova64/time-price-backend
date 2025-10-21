package thing

import (
	"database/sql"
	"log"
	"main/entity"
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

	query := `INSERT INTO thing(name, pay_date, pay_price, sale_date, sale_price) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := db.QueryRow(
		query,
		thing.Name,
		thing.PayDate,
		thing.PayPrice,
		thing.UserID,
		saleDate,
		salePrice,
	).Scan(&thing.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

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
