package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/controllers"
	"main/db"
	"net/http"
)

func main() {
	// Загрузка .env файла
	err := godotenv.Load(".env")
	if err != nil {

		log.Fatal("Main: Ошибка загрузки .env файла.")
	}

	db := db.ConnectDB()
	app := controllers.NewApp(db)

	http.HandleFunc("/", app.HomeController)

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Main: Ошибка сервера: ", err)
	} else {
		log.Println("Main: Сервер запущен на 80-м порту.")
	}
}
