package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"main/controllers"
	"main/db"
	"main/middleware"
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
	router := mux.NewRouter()

	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.AuthMiddleware(db))

	router.HandleFunc("/auth", app.AuthHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/admin/user", app.AdminUserController).Methods("POST")
	router.HandleFunc("/", app.HomeController).Methods("GET")
	router.HandleFunc("/admin/thing", app.AdminThingController).Methods("POST")
	router.HandleFunc("/admin/thing/{id:[0-9]+}", app.AdminThingUpdateController).Methods("PUT")
	router.HandleFunc("/admin/expense", app.AdminExpenseController).Methods("POST")

	err = http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal("Main: Ошибка сервера: ", err)
	} else {
		log.Println("Main: Сервер запущен на 80-м порту.")
	}
}
