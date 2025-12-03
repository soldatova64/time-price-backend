package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"main/entity"
	"main/helpers"
	"main/models"
	"main/models/requests"
	"main/models/responses"
	"main/repositories/expense"
	"net/http"
	"strconv"
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

func (app *App) AdminExpenseUpdateController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_expense_update"}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "id",
					Message: "Неверный формат ID",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

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
		ID:          id,
		ThingID:     req.ThingID,
		Sum:         req.Sum,
		Description: req.Description,
		ExpenseDate: req.ExpenseDate,
	}

	expenseRepo := expense.NewRepository(app.db)
	updatedExpense, err := expenseRepo.Update(&expenseEntity)
	if err != nil {
		log.Printf("Database error: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "database",
					Message: "Не удалось обновить расход в БД",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	successResponse := responses.AdminExpenseResponse{
		Meta: meta,
		Data: *updatedExpense,
	}

	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(successResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}

func (app *App) AdminExpenseDeleteController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_expense_delete"}
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "id",
					Message: "Неверный формат ID",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	expenseRepo := expense.NewRepository(app.db)
	err = expenseRepo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse := responses.ErrorResponse{
				Meta: meta,
				Errors: []responses.Error{
					{
						Field:   "id",
						Message: "Расход с указанным ID не найден",
					},
				},
			}
			writer.WriteHeader(http.StatusNotFound)
			json.NewEncoder(writer).Encode(errorResponse)
			return
		}

		log.Printf("Database error: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "database",
					Message: "Не удалось удалить расход",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	successResponse := responses.AdminExpenseResponse{
		Meta: meta,
	}

	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(successResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}
