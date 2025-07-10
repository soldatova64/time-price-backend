package controllers

import (
	"encoding/json"
	"main/entity"
	"main/models"
	"main/models/responses"
	"main/repositories/auth_token"
	"main/repositories/expense"
	"main/repositories/thing"
	"math"
	"net/http"
	"time"
)

func (app *App) HomeController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "home"}

	token := request.Header.Get("Authorization")
	authToken, err := auth_token.NewRepository(app.db).FindByToken(token)
	if err != nil || authToken.UserID == 0 {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "auth",
					Message: "Необходима авторизация",
				},
			},
		}
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	things, err := thing.NewRepository(app.db).FindAllByUserID(authToken.UserID)
	if err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "database",
					Message: "Ошибка при получении данных",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
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

	for key := range things {
		thingExpenses := expensesByThingID[things[key].ID]
		if thingExpenses == nil {
			thingExpenses = []entity.Expense{}
		}

		var endDate time.Time

		if things[key].SaleDate.Valid {
			endDate = things[key].SaleDate.Time
		} else {
			endDate = dayNow
		}
		things[key].Days = int(endDate.Sub(things[key].PayDate).Hours()/24) + 1

		price := things[key].PayPrice
		if things[key].SalePrice.Valid {
			price -= int(things[key].SalePrice.Int64)
		}

		for _, expense := range thingExpenses {
			price += expense.Sum
		}

		things[key].PayDay = math.Round(float64(price) / float64(things[key].Days))
		things[key].Expense = thingExpenses

	}

	response := responses.HomeResponse{
		Meta: models.Meta{Action: "home"},
		Data: things,
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Ошибка формирования JSON.", http.StatusInternalServerError)
	}

}
