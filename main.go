package main

import (
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()

	router.Use(CorsMiddleware)
	router.HandleFunc("/", app.HomeController)
	router.HandleFunc("/admin/thing", app.AdminThingController).Methods("POST")
	router.HandleFunc("/admin/thing/{id:[0-9]+}", app.AdminThingUpdateController).Methods("PUT")
	router.HandleFunc("/admin/expense", app.AdminExpenseController).Methods("POST")
	router.HandleFunc("/auth", app.AuthHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/admin/user", app.AdminUserController).Methods("POST")

	err = http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal("Main: Ошибка сервера: ", err)
	} else {
		log.Println("Main: Сервер запущен на 80-м порту.")
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
