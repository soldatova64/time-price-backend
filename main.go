package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"main/db"
	"main/handlers"
	"net/http"
)

func main() {
	// Загрузка .env файла
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDB()

	r := chi.NewRouter()

	r.Get("/", handlers.PageHome)
	r.Get("/things", handlers.GetThing)
	r.Post("/things", handlers.PostThing)
	r.Get("/things/{id}", handlers.GetThingByID)
	r.Delete("/things/{id}", handlers.DeleteThing)

	log.Println("Сервер запущен на http://localhost:80")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Ошибка сервера: ", err)
	}
}
