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
	"main/repositories/thing"
	"net/http"
	"strconv"
)

func (app *App) AdminThingController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_thing"}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	userID, ok := request.Context().Value("user_id").(int)
	if !ok {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "auth",
					Message: "Требуется авторизация2",
				},
			},
		}
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	var req requests.ThingRequest
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

	if helpers.IsFutureDate(req.PayDate) {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "pay_date",
					Message: "Дата не может быть в будущем",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	if req.SaleDate.Valid && helpers.IsFutureDate(req.SaleDate.Time) {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "sale_date",
					Message: "Дата продажи не может быть в будущем",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	thingEntity := entity.Thing{
		Name:      req.Name,
		PayDate:   req.PayDate,
		PayPrice:  req.PayPrice,
		SaleDate:  req.SaleDate,
		SalePrice: req.SalePrice,
		UserID:    userID,
	}

	thingRepo := thing.NewRepository(app.db)
	createdThing, err := thingRepo.Add(&thingEntity)
	if err != nil {
		log.Printf("Database error: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "database",
					Message: "Не удалось добавить вещь в БД",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	successResponse := responses.AdminThingResponse{
		Meta: meta,
		Data: *createdThing,
	}
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(successResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}

func (app *App) AdminThingUpdateController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_thing_update"}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	userID, ok := request.Context().Value("user_id").(int)
	if !ok {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "auth",
					Message: "Требуется авторизация3",
				},
			},
		}
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

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

	thingRepo := thing.NewRepository(app.db)
	thingEntity, err := thingRepo.Find(id, userID)
	if err != nil {
		log.Printf("Database error: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "database",
					Message: "Не удалось найти вещь в БД",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	var req requests.ThingRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "request",
					Message: "Недопустимый формат запроса5",
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

	if helpers.IsFutureDate(req.PayDate) {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "pay_date",
					Message: "Дата не может быть в будущем",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	if req.SaleDate.Valid && helpers.IsFutureDate(req.SaleDate.Time) {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "sale_date",
					Message: "Дата продажи не может быть в будущем",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	thingEntity.Name = req.Name
	thingEntity.PayDate = req.PayDate
	thingEntity.PayPrice = req.PayPrice
	thingEntity.SaleDate = req.SaleDate
	thingEntity.SalePrice = req.SalePrice

	updatedThing, err := thingRepo.Update(thingEntity)
	if err != nil {
		log.Printf("Database error: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "database",
					Message: "Не удалось обновить вещь в БД",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	successResponse := responses.AdminThingResponse{
		Meta: meta,
		Data: updatedThing,
	}
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(successResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}

func (app *App) AdminThingDeleteController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_thing_delete"}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	userID, ok := request.Context().Value("user_id").(int)
	if !ok {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "auth",
					Message: "Требуется авторизация",
				},
			},
		}
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

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

	thingRepo := thing.NewRepository(app.db)
	err = thingRepo.Delete(id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse := responses.ErrorResponse{
				Meta: meta,
				Errors: []responses.Error{
					{
						Field:   "id",
						Message: "Вещь не найдена",
					},
				},
			}
			writer.WriteHeader(http.StatusNotFound)
			json.NewEncoder(writer).Encode(errorResponse)
		} else {
			log.Printf("Database error: %v", err)
			errorResponse := responses.ErrorResponse{
				Meta: meta,
				Errors: []responses.Error{
					{
						Field:   "database",
						Message: "Не удалось удалить вещь",
					},
				},
			}
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(errorResponse)
		}
		return
	}

	successResponse := responses.AdminThingResponse{
		Meta: meta,
	}
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(successResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}
