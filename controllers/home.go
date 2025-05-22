package controllers

import (
	"encoding/json"
	"main/models"
	"main/models/responses"
	"main/repositories"
	"math"
	"net/http"
	"time"
)

func (app *App) HomeController(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	collection, err := repositories.FindAll(app.db)

	if err != nil {
		http.Error(writer, "Ошибка базы данных.", http.StatusInternalServerError)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	dayNow := time.Now()
	dayNow = time.Date(dayNow.Year(), dayNow.Month(), dayNow.Day(), 0, 0, 0, 0, time.UTC)

	for key := range collection {
		var endDate time.Time

		if collection[key].SaleDate.Valid {
			endDate = collection[key].SaleDate.Time
		} else {
			endDate = dayNow
		}

		collection[key].Days = int(endDate.Sub(collection[key].PayDate).Hours()/24) + 1

		var prise int
		if collection[key].SalePrice.Valid {
			prise = collection[key].PayPrice - int(collection[key].SalePrice.Int64)
		} else {
			prise = collection[key].PayPrice
		}

		collection[key].PayDay = float64(prise) / float64(collection[key].Days)
		collection[key].PayDay = math.Round(collection[key].PayDay)
	}

	response := responses.HomeResponse{
		Meta: models.Meta{Action: "home"},
		Data: collection,
	}

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Ошибка формирования JSON.", http.StatusInternalServerError)
	}
}
