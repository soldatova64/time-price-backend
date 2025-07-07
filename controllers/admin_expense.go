package controllers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"main/entity"
	"main/helpers"
	"main/models"
	"main/models/requests"
	"main/models/responses"
	"main/repositories/expense"
	"net/http"
)

func (app *App) AdminExpenseController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_expense"}
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	var req requests.ExpenseRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "request",
					Message: "Недопустимый формат запроса",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		errors := helpers.ParseValidationErrors(err)
		errorResponse := responses.ErrorResponse{
			Meta:   meta,
			Errors: errors,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	if helpers.IsFutureDate(req.ExpenseDate) {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "expense_date",
					Message: "Дата не может быть в будущем",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	expenseEntity := entity.Expense{
		ThingID:     req.ThingID,
		Sum:         req.Sum,
		Description: req.Description,
		ExpenseDate: req.ExpenseDate,
	}

	expenseRepo := expense.NewRepository(app.db)
	createdExpense, err := expenseRepo.Add(&expenseEntity)
	if err != nil {
		log.Printf("Database error: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "database",
					Message: "Не удалось добавить расход в БД",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	successResponse := responses.AdminExpenseResponse{
		Meta: meta,
		Data: *createdExpense,
	}

	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(successResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}
