package controllers

import (
	"encoding/json"
	"main/models"
	"main/models/responses"
	"main/repositories"
	"net/http"
)

func (app *App) HomeController(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	collection, err := repositories.FindAll(app.db)

	if err != nil {
		http.Error(writer, "Ошибка базы данных.", http.StatusInternalServerError)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.HomeResponse{
		Meta: models.Meta{Action: "home"},
		Data: collection,
	}

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Ошибка формирования JSON.", http.StatusInternalServerError)
	}
}
