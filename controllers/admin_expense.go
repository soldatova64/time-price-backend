package controllers

import (
	"encoding/json"
	"log"
	"main/entity"
	"main/models"
	"main/models/responses"
	"main/repositories/expense"
	"net/http"
	"time"
)

func (app *App) AdminExpenseController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_expense"}
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	var expenseEntity entity.Expense

	if err := json.NewDecoder(request.Body).Decode(&expenseEntity); err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "json",
					Message: "Недопустимый формат JSON",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	var errors []responses.Error

	if expenseEntity.Sum <= 0 {
		errors = append(errors, responses.Error{
			Field:   "sum",
			Message: "Сумма расхода должна быть положительной",
		})
	}

	if expenseEntity.Description == "" {
		errors = append(errors, responses.Error{
			Field:   "description",
			Message: "Описание расхода должно быть заполнено",
		})
	}

	if expenseEntity.ExpenseDate.IsZero() {
		errors = append(errors, responses.Error{
			Field:   "expense_date",
			Message: "Необходимо указать дату расхода",
		})
	} else if expenseEntity.ExpenseDate.After(time.Now()) {
		errors = append(errors, responses.Error{
			Field:   "expense_date",
			Message: "Дата не может быть в будущем",
		})
	}

	if len(errors) > 0 {
		errorResponse := responses.ErrorResponse{
			Meta:   meta,
			Errors: errors,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
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
