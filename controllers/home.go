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

	type ThingResponse struct {
		entity.Thing
		Expenses []entity.Expense `json:"expenses"`
	}

	var responseThings []ThingResponse

	dayNow := time.Now()
	dayNow = time.Date(dayNow.Year(), dayNow.Month(), dayNow.Day(), 0, 0, 0, 0, time.UTC)

	for _, key := range things {
		thingExpenses := expensesByThingID[key.ID]
		if thingExpenses == nil {
			thingExpenses = []entity.Expense{}
		}

		var endDate time.Time

		if key.SaleDate.Valid {
			endDate = key.SaleDate.Time
		} else {
			endDate = dayNow
		}

		key.Days = int(endDate.Sub(key.PayDate).Hours()/24) + 1

		var prise int
		if key.SalePrice.Valid {
			prise = key.PayPrice - int(key.SalePrice.Int64)
		} else {
			prise = key.PayPrice
		}

		for _, expense := range thingExpenses {
			prise += expense.Sum
		}

		key.PayDay = float64(prise) / float64(key.Days)
		key.PayDay = math.Round(key.PayDay)

		thingCopy := entity.Thing{
			ID:        key.ID,
			Name:      key.Name,
			PayDate:   key.PayDate,
			PayPrice:  key.PayPrice,
			SaleDate:  key.SaleDate,
			SalePrice: key.SalePrice,
			Days:      key.Days,
			PayDay:    key.PayDay,
		}

		responseThings = append(responseThings, ThingResponse{
			Thing:    thingCopy,
			Expenses: thingExpenses,
		})
	}

	response := responses.HomeResponse{
		Meta: models.Meta{Action: "home"},
		Data: struct {
			Things []ThingResponse `json:"things"`
		}{
			Things: responseThings,
		},
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Ошибка формирования JSON.", http.StatusInternalServerError)
	}

}
