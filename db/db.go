package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func ConnectDB() (db *sql.DB) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
		log.Println("DB: DB_HOST not set, using default 'localhost'.")
	}

	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")

	log.Println("DB: Migrate begin.")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	if err == migrate.ErrNoChange {
		log.Println("DB: Migrate no change.")
	}

	log.Println("DB: Migrate done.")

	log.Println(fmt.Sprintf("DB: Connection to %s:%s/%s", host, port, dbname))

	dsn = fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", host, user, dbname, pass, port)

	db, err = sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("DB: Failed to connect to database.\n", err)
		os.Exit(1)
	}

	log.Println("DB: Connected.")

	return db
}
