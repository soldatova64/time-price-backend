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

	// Публичные роуты (без аутентификации)
	publicRouter := router.PathPrefix("").Subrouter()
	publicRouter.HandleFunc("/auth", app.AuthHandler).Methods("POST", "OPTIONS")
	publicRouter.HandleFunc("/admin/user", app.AdminUserController).Methods("POST")

	// Защищенные роуты (требуют аутентификации)
	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(db))
	protectedRouter.HandleFunc("/", app.HomeController).Methods("GET")
	protectedRouter.HandleFunc("/admin/thing", app.AdminThingController).Methods("POST")
	protectedRouter.HandleFunc("/admin/thing/{id:[0-9]+}", app.AdminThingUpdateController).Methods("PUT")
	protectedRouter.HandleFunc("/admin/expense", app.AdminExpenseController).Methods("POST")
	err = http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal("Main: Ошибка сервера: ", err)
	} else {
		log.Println("Main: Сервер запущен на 80-м порту.")
	}
}
