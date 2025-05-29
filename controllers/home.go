package controllers

import (
	"encoding/json"
	"main/entity"
	"main/models"
	"main/models/responses"
	"main/repositories/expense"
	"main/repositories/thing"
	"math"
	"net/http"
	"time"
)

func (app *App) HomeController(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	things, err := thing.NewRepository(app.db).FindAll()
	if err != nil {
		http.Error(writer, "Ошибка базы данных.", http.StatusInternalServerError)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	expenses, err := expense.NewRepository(app.db).FindAll()
	if err != nil {
		http.Error(writer, "Ошибка базы данных.", http.StatusInternalServerError)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	expensesByThingID := make(map[int][]entity.Expense)
	for _, e := range expenses {
		expensesByThingID[e.ThingID] = append(expensesByThingID[e.ThingID], e)
	}

	dayNow := time.Now()
	dayNow = time.Date(dayNow.Year(), dayNow.Month(), dayNow.Day(), 0, 0, 0, 0, time.UTC)

	for i := range things {
		thingExpenses := expensesByThingID[things[i].ID]
		if thingExpenses == nil {
			thingExpenses = []entity.Expense{}
		}

		var endDate time.Time

		if things[i].SaleDate.Valid {
			endDate = things[i].SaleDate.Time
		} else {
			endDate = dayNow
		}
		things[i].Days = int(endDate.Sub(things[i].PayDate).Hours()/24) + 1

		price := things[i].PayPrice
		if things[i].SalePrice.Valid {
			price -= int(things[i].SalePrice.Int64)
		}

		for _, expense := range thingExpenses {
			price += expense.Sum
		}

		things[i].PayDay = math.Round(float64(price) / float64(things[i].Days))
		things[i].Expense = thingExpenses

	}

	response := responses.HomeResponse{
		Meta: models.Meta{Action: "home"},
		Data: things,
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Ошибка формирования JSON.", http.StatusInternalServerError)
	}

}
