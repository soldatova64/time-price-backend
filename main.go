package main

import (
	"encoding/json"
	"log"
	"main/db"
	"net/http"
)

type Meta struct {
	Action string `json:"action"`
}

type Response struct {
	Meta Meta                   `json:"meta"`
	Data map[string]interface{} `json:"data"`
}

func main() {
	db.ConnectDB()

	http.HandleFunc("/", PageHome)
	log.Println("Сервер запущен на http://localhost:80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Ошибка сервера: ", err)
	}
}

func PageHome(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	response := Response{
		Meta: Meta{Action: "home"},
		Data: make(map[string]interface{}),
	}

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}
